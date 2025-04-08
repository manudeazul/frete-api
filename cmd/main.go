package main

import (
	"log"
	"superfrete-api/handler"
	"superfrete-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar .env")
	}

	server := gin.Default()
	dbConnection, err := repository.ConnectDB()
	if err != nil {
		log.Fatalf("Erro ao iniciar Gin: %v", err)
	}

	quoteRepo := repository.NewQuoteRepository(dbConnection)
	quoteHandler := handler.NewQuoteHandler(quoteRepo)

	server.GET("/metrics", quoteHandler.GetLastQuote)
	server.POST("/quote", quoteHandler.PostQuote)

	err = server.Run(":8000")
	if err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
