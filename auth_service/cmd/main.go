package main

import (
	"auth_service/configs"
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

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var config configs.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	db, dbInterface, err := repository.ConnectToDb(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	}
	defer dbInterface.Close(db)
	rdb, redisInterface, err := repository.ConnectToRedis(config)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return
	}
	defer redisInterface.Close(rdb)

	reposdb := repository.NewAuthRepository(db)

	service := service.NewService(reposdb)
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
