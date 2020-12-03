package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"example.com/auth-service-go/api/handler"
	"example.com/auth-service-go/config"
	"example.com/auth-service-go/internal/infrastructure/database"
	mongo "example.com/auth-service-go/internal/repository/token/mongo"
	"github.com/go-chi/chi"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	log.Println("Starting the server")

	cfg := config.New()
	ctx := context.Background()
	mongoDB, ctx := database.NewMongoClient(ctx, cfg)
	defer mongoDB.Disconnect(ctx)

	tokenMongoRepo := mongo.NewTokenRepository(mongoDB, "tokens")
	router := chi.NewRouter()

	handler := handler.New(ctx, router)
	handler.InitAuthRoutes(tokenMongoRepo)

	s := http.Server{
		Addr:         ":" + cfg.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      handler.Router,
	}

	log.Println("Server is running")
	return s.ListenAndServe()
}
