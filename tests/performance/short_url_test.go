package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"knands42/url-shortener/internal/cache"
	"knands42/url-shortener/internal/database"
	"knands42/url-shortener/internal/database/repo"
	handler "knands42/url-shortener/internal/handlers"
	"knands42/url-shortener/internal/otel"
	"knands42/url-shortener/internal/server"
	"knands42/url-shortener/internal/utils"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"golang.org/x/exp/rand"
)

var r = chi.NewRouter()
var ctx = context.Background()
var config = utils.NewConfig("test")
var dbConfig = database.NewDBConfig(
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
var dbConnection, _ = dbConfig.Connect(ctx)
var repository = repo.New(dbConnection)
var caching = cache.NewRedisClient(config)
var exporter, _ = stdouttrace.New()
var tracing = otel.NewOpenTelemetry(exporter).GetTracer()
var handlers = handler.NewHandler(repository, caching, tracing)
var testServer = server.NewServer(r, handlers, tracing)

func Benchmark_short_url(b *testing.B) {
	ts := httptest.NewServer(testServer.Router)
	creationAmount := 1000
	getAmount := 10000
	var url string
	var urls []string

	b.Run("CreateEntry", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for j := 0; j < creationAmount; j++ {
				url = generateEntry(b, ts)
				urls = append(urls, url)
			}
		}
	})

	b.Run("GetEntry", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			semaphore := make(chan struct{}, 25)
			var wg sync.WaitGroup

			for j := 0; j < getAmount; j++ {
				wg.Add(1)
				semaphore <- struct{}{}

				go func() {
					defer wg.Done()
					defer func() { <-semaphore }()

					randomUrl := urls[rand.Intn(len(urls))]
					getEntry(b, ts, randomUrl)
				}()
			}

			wg.Wait()
		}
	})

}

func generateEntry(b *testing.B, ts *httptest.Server) string {
	// Define the request payload
	url := "https://" + generateRandomString(10) + ".com"
	payload := map[string]string{"input": url}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		b.Fatalf("Failed to marshal payload: %v", err)
	}

	// Send a POST request to the /api/v1/shorten endpoint
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/shorten", bytes.NewBuffer(payloadBytes))
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		b.Fatalf("Failed to send POST request: %v", err)
	}

	return url
}

func getEntry(b *testing.B, ts *httptest.Server, url string) {
	// Send a POST request to the /api/v1/shorten endpoint
	req, _ := http.NewRequest("GET", ts.URL+"/api/v1/url?url="+url+"&type=original_url", nil)
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		b.Fatalf("Failed to send POST request: %v", err)
	}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
