package main

import (
	"context"
	"fmt"
	"knands42/url-shortener/internal/cache"
	"knands42/url-shortener/internal/database"
	"knands42/url-shortener/internal/database/repo"
	handler "knands42/url-shortener/internal/handlers"
	"knands42/url-shortener/internal/otel"
	"knands42/url-shortener/internal/server"
	"knands42/url-shortener/internal/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "knands42/url-shortener/docs"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

//	@title			Url Shortener API
//	@version		1.0
//	@description	This is a sample server for a URL Shortener API.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:3333
// @BasePath	/api/v1
// @schemes	http
// @produce	json
// @consumes	json
// @failure	500 {object} utils.InternalServerErrorResponse
func main() {
	ctx := context.Background()
	r := chi.NewRouter()
	config := utils.NewConfig("dev")
	jaeger, tracing := InitTracing(ctx, config)
	defer func() { _ = jaeger.Shutdown(ctx) }()
	dbConnection := InitDatabase(ctx, config)
	defer dbConnection.Close()

	// Intialized the components
	repo := repo.New(dbConnection)
	cache := cache.NewRedisClient(config)
	handlers := handler.NewHandler(repo, cache, tracing)
	server.NewServer(r, handlers, tracing)

	gracefulShutdown(dbConnection, jaeger, cache)

	err := http.ListenAndServe(":3333", r)
	if err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(conn *pgxpool.Pool, jaeger sdkTrace.SpanExporter, cache *redis.Client) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+C
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)
	go func() {
		// Wait for a signal
		<-signalChan

		fmt.Println("Shutting down the server...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		select {
		case <-time.After(30 * time.Second):
			conn.Close()
			jaeger.Shutdown(ctx)
			cache.Close()
		case <-ctx.Done():
			fmt.Println("Graceful shutdown")
		}
	}()
}

func InitTracing(ctx context.Context, config *utils.Config) (sdkTrace.SpanExporter, trace.Tracer) {
	jaeger, err := otel.NewJaegerExporter(ctx, config.JaegerExporterEndpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize Jaeger: %v\n", err)
		os.Exit(1)
	}

	return jaeger, otel.NewOpenTelemetry(jaeger).GetTracer()
}

func InitDatabase(ctx context.Context, config *utils.Config) *pgxpool.Pool {
	dbConfig := database.NewDBConfig(
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
		config.SSLMode,
		config.TimeZone,
		config.MinConns,
		config.MaxConns,
	)
	dbConnection, err := dbConfig.Connect(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return dbConnection
}
