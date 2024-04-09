package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var err error

type Books struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

func main() {
	db, err = gorm.Open("sqlite3", "./kitaplar.db")

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.AutoMigrate(&Books{})

	r := gin.Default()

	r.POST("/allbooks", CreateBook)
	r.GET("/allbooks/", GetBooks)
	r.GET("/allbooks/:id", GetBook)
	r.GET("/allbooks/authors/:author", GetAuthor)
	r.PUT("/allbooks/:id", UpdateBook)
	r.DELETE("/allbooks/:id", DeleteBook)

	r.Run(":8080")
}
func CreateBook(c *gin.Context) {
	var book Books
	c.BindJSON(&book)
	db.Create(&book)
	c.JSON(200, book)
}

func GetBooks(c *gin.Context) {
	var books []Books

	if err := db.Find(&books).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, books)
	}
}

func GetBook(c *gin.Context) {
	id := c.Params.ByName("id")
	var book Books
	if err := db.Where("id = ?", id).First(&book).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, book)
	}
}

func GetAuthor(c *gin.Context) {
	author := c.Params.ByName("author")
	var book []Books

	if err := db.Where("author = ?", author).Find(&book).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}
	if len(book) > 0 {
		z := func(a []Books) []string {
			e := make([]string, 0)
			for _, v := range a {
				e = append(e, v.Name)
			}
			return e
		}
		n := z(book)
		n = append(n, book[0].Author)
		c.JSON(200, n)
	} else {
		c.JSON(200, book)
	}
}

func UpdateBook(c *gin.Context) {
	var book Books
	id := c.Params.ByName("id")
	if err := db.Where("id = ?", id).First(&book).Error; err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	}

	c.BindJSON(&book)
	db.Save(&book)
	c.JSON(200, book)
}

func DeleteBook(c *gin.Context) {
	id := c.Params.ByName("id")
	var book Books
	d := db.Where("id = ?", id).Delete(&book)
	fmt.Println(d)
	c.JSON(200, gin.H{"id #" + id: "deleted"})
}
