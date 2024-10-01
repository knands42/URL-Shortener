package main

import (
	"context"
	"fmt"
	"knands42/url-shortener/internal/database"
	handler "knands42/url-shortener/internal/handlers"
	"knands42/url-shortener/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	// TODO: Load env environments

	// Initialize the database
	dbConfig := database.NewDBConfig()
	_, err := dbConfig.Connect()
	if err != nil {
		panic(err)
	}

	// Initialize the handlers
	handlers := handler.NewHandler()

	// Initialize the server
	server.NewServer(r, handlers)
	gracefulShutdown()

	err = http.ListenAndServe(":3333", r)
	if err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown() {
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
			// TODO close connection
		case <-ctx.Done():
			fmt.Println("Graceful shutdown")
		}
	}()
}
