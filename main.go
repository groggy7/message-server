package main

import (
	"context"
	"message-server/internal/controller/router"
	"message-server/internal/repository"
	"message-server/internal/usecases"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		panic("DB_URL not set in .env")
	}

	firebaseCredentials := os.Getenv("FIREBASE_CREDENTIALS")
	if firebaseCredentials == "" {
		panic("FIREBASE_CREDENTIALS not set in .env")
	}

	firebaseBucket := os.Getenv("FIREBASE_BUCKET")
	if firebaseBucket == "" {
		panic("FIREBASE_BUCKET not set in .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	poolConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		panic(err)
	}

	roomRepository := repository.NewRoomRepository(pool)
	authRepository := repository.NewAuthRepository(pool)
	listingRepository := repository.NewListingRepository(pool)
	fileRepository := repository.NewFileRepository(firebaseCredentials, firebaseBucket)

	roomUseCase := usecases.NewRoomUseCase(roomRepository)
	authUseCase := usecases.NewAuthUseCase(authRepository)
	listingUseCase := usecases.NewListingUseCase(listingRepository)
	fileUseCase := usecases.NewFileUseCase(fileRepository)

	router := router.NewRouter(roomUseCase, authUseCase, listingUseCase, fileUseCase)
	router.Run(":" + port)
}
