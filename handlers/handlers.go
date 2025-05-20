package handlers

import (
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "book-api/models"
    "book-api/storage"
)

// GetBooks lists all books
func GetBooks(c *gin.Context) {
    books := storage.GetBooks()
    c.JSON(http.StatusOK, books)
}

// CreateBook adds a new book
func CreateBook(c *gin.Context) {
    var newBook models.Book
    if err := c.ShouldBindJSON(&newBook); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    createdBook := storage.CreateBook(newBook)
    if err := storage.SaveBooks(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save books"})
        return
    }
    c.JSON(http.StatusCreated, createdBook)
}

// GetBookByID retrieves a single book
func GetBookByID(c *gin.Context) {
    id := c.Param("id")
    book, found := storage.GetBookByID(id)
    if !found {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }
    c.JSON(http.StatusOK, book)
}

// UpdateBook updates a book
func UpdateBook(c *gin.Context) {
    id := c.Param("id")
    var updatedBook models.Book
    if err := c.ShouldBindJSON(&updatedBook); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    if !storage.UpdateBook(id, updatedBook) {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }
    if err := storage.SaveBooks(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save books"})
        return
    }
    c.JSON(http.StatusOK, updatedBook)
}

// DeleteBook deletes a book
func DeleteBook(c *gin.Context) {
    id := c.Param("id")
    if !storage.DeleteBook(id) {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }
    if err := storage.SaveBooks(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save books"})
        return
    }
    c.Status(http.StatusNoContent)
}

// SearchBooks searches books by keyword with concurrency
func SearchBooks(c *gin.Context) {
    keyword := c.Query("q")
    if keyword == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search keyword"})
        return
    }

    books := storage.GetBooks()
    if len(books) == 0 {
        c.JSON(http.StatusOK, []models.Book{})
        return
    }

    numGoroutines := 4
    if len(books) < numGoroutines {
        numGoroutines = len(books) // Avoid unnecessary goroutines for small lists
    }
    chunkSize := (len(books) + numGoroutines - 1) / numGoroutines
    resultsChan := make(chan []models.Book, numGoroutines)

    for i := 0; i < numGoroutines; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if end > len(books) {
            end = len(books)
        }
        go searchChunk(books[start:end], keyword, resultsChan)
    }

    var results []models.Book
    for i := 0; i < numGoroutines; i++ {
        chunkResults := <-resultsChan
        results = append(results, chunkResults...)
    }
    c.JSON(http.StatusOK, results)
}

// searchChunk searches a subset of books
func searchChunk(books []models.Book, keyword string, resultsChan chan<- []models.Book) {
    var chunkResults []models.Book
    for _, book := range books {
        if strings.Contains(strings.ToLower(book.Title), strings.ToLower(keyword)) ||
            strings.Contains(strings.ToLower(book.Description), strings.ToLower(keyword)) {
            chunkResults = append(chunkResults, book)
        }
    }
    resultsChan <- chunkResults
}