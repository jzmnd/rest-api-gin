package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// album represents data about a record album.
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dbpool := dbConnect(ctx)
	defer dbpool.Close()
	log.Println("Connected to database")

	router := gin.Default()
	router.GET("/ping", handlePing)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("0.0.0.0:8080")
}

// dbConnect creates the database connection.
func dbConnect(c context.Context) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig("postgres://")
	if err != nil {
		log.Fatalf("Unable to parse config: %v\n", err)
	}

	// Update config using database parameters from the environment.
	config.ConnConfig.Host = os.Getenv("DB_HOST")
	config.ConnConfig.User = os.Getenv("DB_USER")
	config.ConnConfig.Password = os.Getenv("DB_PASSWORD")
	config.ConnConfig.Database = os.Getenv("DB_NAME")

	// Create the database connection pool.
	dbpool, err := pgxpool.NewWithConfig(c, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	return dbpool
}

// handlePing responds with a healthcheck message.
func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum Album

	// Call BindJSON to bind the received JSON to newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
