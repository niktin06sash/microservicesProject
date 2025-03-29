package main

import (
	"auth_service/internal/api"
	"auth_service/internal/repository"
	"auth_service/internal/server"
	"auth_service/internal/service"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func main() {

	configPath := "../configs/config.yml"
	repository.LoadConfig(configPath)

	dbInterface := repository.DBObject{}
	db, err := repository.ConnectToDb(dbInterface)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	defer dbInterface.Close(db)

	repos := repository.NewAuthRepository(db)
	service := service.NewService(repos)
	handlers := api.NewHandler(service)
	srv := &server.Server{}

	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting server on port: %s\n", port)
	serverError := make(chan error, 1)
	go func() {

		if err := srv.Run(port, handlers.InitRoutes()); err != nil {
			serverError <- fmt.Errorf("server run failed: %w", err)
			return
		}
		close(serverError)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case sig := <-quit:
		log.Printf("Service shutting down with signal: %v", sig)
	case err := <-serverError:
		log.Fatalf("Service startup failed: %v", err)
	}

	shutdownTimeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	log.Println("Service is shutting down...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Service has shutted down successfully")

}
