package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "Golang", Author: "Mic", Quantity: 10},
	{ID: "2", Title: "Java", Author: "Tom", Quantity: 10},
	{ID: "3", Title: "SQL", Author: "Smith", Quantity: 10},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request Body"})
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {

	for i, book := range books {
		if book.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func checkOutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
	}

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not avaiable"})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
	}

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/newbook", createBook)
	router.GET("book/:id", bookById)
	router.PATCH("/checkout", checkOutBook)
	router.PATCH("/return", returnBook)

	router.Run("localhost:8080")
}
