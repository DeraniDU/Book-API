package storage

import (
    "encoding/json"
    "os"
    "github.com/google/uuid"
    "book-api/models"
)

const fileName = "books.json"

var books []models.Book

// LoadBooks reads books from the JSON file
func LoadBooks() error {
    file, err := os.Open(fileName)
    if err != nil {
        if os.IsNotExist(err) {
            books = []models.Book{} // Initialize empty slice if file doesn’t exist
            return nil
        }
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    return decoder.Decode(&books)
}

// SaveBooks writes books to the JSON file
func SaveBooks() error {
    file, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    return encoder.Encode(books)
}

// GetBooks returns all books
func GetBooks() []models.Book {
    return books
}

// GetBookByID finds a book by ID
func GetBookByID(id string) (models.Book, bool) {
    for _, book := range books {
        if book.BookId == id {
            return book, true
        }
    }
    return models.Book{}, false
}

// CreateBook adds a new book
func CreateBook(book models.Book) models.Book {
    book.BookId = uuid.New().String() // Generate unique ID
    books = append(books, book)
    return book
}

// UpdateBook updates an existing book
func UpdateBook(id string, updatedBook models.Book) bool {
    for i, book := range books {
        if book.BookId == id {
            updatedBook.BookId = id
            books[i] = updatedBook
            return true
        }
    }
    return false
}

// DeleteBook removes a book by ID
func DeleteBook(id string) bool {
    for i, book := range books {
        if book.BookId == id {
            books = append(books[:i], books[i+1:]...)
            return true
        }
    }
    return false
}