package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Darrell-Devana/be-note-app/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Note struct {
	Title       string          `json:"title"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
	LastVisited string          `json:"lastVisited"`
	TextContent string          `json:"textContent"`
	Content     json.RawMessage `json:"content"`
	ID          uuid.UUID       `json:"id"`
	IsFavorite  bool            `json:"isFavorite"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	database.ConnectDatabase()

	r.POST("/add", addNote)
	r.GET("/list", listNote)
	r.GET("/notes/:id", openNote)
	r.POST("/update", updateNote)
	r.POST("/delete/:id", deleteNote)
	r.POST("/favorite/:id", favoriteNote)
	r.Run(":8008")
}

func addNote(c *gin.Context) {
	note := Note{}
	if err := c.ShouldBindJSON(&note); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Title and content is required"})
		return
	}

	var noteID string
	if errSql := database.DB.QueryRow(
		"INSERT INTO notes(content, title, is_favorite, text_content) VALUES ($1, $2, $3, $4) RETURNING id",
		note.Content, note.Title, note.IsFavorite, note.TextContent,
	).Scan(&noteID); errSql != nil {
		fmt.Println(errSql)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database operation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Note created successfully", "id": noteID})
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
		err := rows.Scan(&note.ID, (*[]byte)(&note.Content), &note.Title, &note.IsFavorite, &note.CreatedAt, &note.UpdatedAt, &note.LastVisited, &note.TextContent)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to retrieve notes"})
			return
		}
		notes = append(notes, note)
	}
	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "List successful", "output": notes})
}

func openNote(c *gin.Context) {
	id := c.Param("id")
	fmt.Printf("%v\n", id)
	rows, err := database.DB.Query("select * from notes where id = $1", id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database operation failed"})
		return
	}

	notes := []Note{}
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, (*[]byte)(&note.Content), &note.Title, &note.IsFavorite, &note.CreatedAt, &note.UpdatedAt, &note.LastVisited, &note.TextContent)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to retrieve notes"})
			return
		}
		notes = append(notes, note)
	}

	if _, err = database.DB.Exec("update notes set last_visited = now() where id = $1", notes[0].ID); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database operation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Note found successfully", "output": notes[0]})
}

func updateNote(c *gin.Context) {
	note := Note{}
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "400", "message": "Invalid request body. Please ensure all required fields are provided and correctly formatted."})
		return
	}

	fmt.Println(note)

	if _, err := database.DB.Exec("update notes set content = $2, title = $3, is_favorite = $4, created_at = $5, updated_at = now(), last_visited = now(), text_content = $6 where id = $1", note.ID, note.Content, note.Title, note.IsFavorite, note.CreatedAt, note.TextContent); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database operation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Note updated successfully", "output": note})
}

func favoriteNote(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)

	rows, err := database.DB.Query("select * from notes where id = $1", id)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database operation failed"})
		return
	}

	notes := []Note{}
	for rows.Next() {
		var note Note
		err := rows.Scan(&note.ID, (*[]byte)(&note.Content), &note.Title, &note.IsFavorite, &note.CreatedAt, &note.UpdatedAt, &note.LastVisited, &note.TextContent)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Failed to retrieve notes"})
			return
		}
		notes = append(notes, note)
	}

	if notes[0].IsFavorite {
		_, err = database.DB.Exec("update notes set is_favorite = $2 where id = $1", id, false)
	} else {
		_, err = database.DB.Exec("update notes set is_favorite = $2 where id = $1", id, true)
	}

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database operation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Note deleted successfully", "deleted": id})
}

func deleteNote(c *gin.Context) {
	id := c.Param("id")
	fmt.Println(id)

	if _, err := database.DB.Exec("delete from notes where id = $1", id); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": "500", "message": "Database operation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": "200", "message": "Note deleted successfully", "deleted": id})
}
