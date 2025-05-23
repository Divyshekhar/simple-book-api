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
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

//c is the all the information about the request {like in hono}

func createBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil { //BindJson takes the json body of the request and decode it to the newBook struct
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book not found")
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		//gin.H is nothing just a map of string map[string] which is the format of json like json({message: "message"}) here {"message":"message"}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing query parameters"})
		return
	}
	book, er := getBookById(id)
	if er != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}

func returnBook(c *gin.Context){
	id, ok := c.GetQuery("id")
	if !ok{
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "query not found"})
	}
	book, err := getBookById(id)
	if err != nil{
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
	}
	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)

}

func main() {
	router := gin.Default()

	router.GET("/books", getBooks)

	router.POST("/create", createBook)

	router.GET("/books/:id", bookById)

	router.PATCH("/checkout", checkoutBook)

	router.PATCH("/return", returnBook)

	router.Run("localhost:5000")

}
