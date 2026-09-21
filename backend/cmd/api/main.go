package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-go/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system environment")
	}

	log := logger.NewLogger(&logger.LoggerOptions{
		AppName: os.Getenv("APP_NAME"),
	})

	handler, cleanup, err := Init(log)
	if err != nil {
		log.WithError(err).Fatal("failed to create server")
	}
	defer cleanup()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Infof("server starting on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Fatal("failed to listen and serve")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := srv.Shutdown(ctx); err != nil {
		cancel()
		log.WithError(err).Fatal("forced shutdown")
		return
	}

	cancel()
	log.Info("server stopped")
}
