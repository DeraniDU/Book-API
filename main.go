package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "book-api/handlers"
    "book-api/storage"
)

func main() {
    // Load books from file at startup
    if err := storage.LoadBooks(); err != nil {
        log.Fatalf("Failed to load books: %v", err)
    }

    // Set up Gin router
    r := gin.Default()

    // Register CRUD endpoints
    r.GET("/books", handlers.GetBooks)
    r.POST("/books", handlers.CreateBook)
    r.GET("/books/:id", handlers.GetBookByID)
    r.PUT("/books/:id", handlers.UpdateBook)
    r.DELETE("/books/:id", handlers.DeleteBook)

    // Register search endpoint
    r.GET("/books/search", handlers.SearchBooks)

    // Start the server
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}