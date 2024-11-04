package main

import (
	"fmt"
	"net/http"

	"github.com/Darrell-Devana/be-note-app/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Note struct {
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	CreatedAt  string    `json:"createdAt"`
	UpdatedAt  string    `json:"updatedAt"`
	ID         uuid.UUID `json:"id"`
	IsFavorite bool      `json:"isFavorite"`
}

func addNote(c *gin.Context) {
	note := Note{}
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Title and content is required"})
		return
	}

	if note.Content == "" || note.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Title and content is required"})
		return
	}

	fmt.Println("Title:", note.Title, "Content:", note.Content)

	_, errSql := database.DB.Exec("insert into notes(content, title, is_favorite) values ($1, $2, $3)", note.Content, note.Title, note.IsFavorite)
	if errSql != nil {
		fmt.Println(errSql)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database operation failed"})
		return
	}

	c.String(http.StatusOK, note.Content)
}

func listNote(c *gin.Context) {
	rows, err := database.DB.Query("select * from notes")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database operation failed"})
		return
	}

	notes := []Note{}
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, &note.Content, &note.Title, &note.IsFavorite, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to retrieve notes"})
			return
		}
		notes = append(notes, note)
	}
	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "List successful", "output": notes})
}

func helloWorld(c *gin.Context) {
	c.String(http.StatusOK, "Hello, world!\n")
}

func main() {
	r := gin.Default()

	r.Use(cors.Default())

	err := godotenv.Load("application.env")
	if err != nil {
		panic(err)
	}
	database.ConnectDatabase()
	r.GET("/hello", helloWorld)
	r.POST("/add", addNote)
	r.GET("/list", listNote)
	r.Run(":8008")
}
