package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jzmnd/rest-api-gin/models"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Env struct {
	Albums interface {
		GetAll(context.Context) ([]models.Album, error)
		GetByID(context.Context, int) (*models.Album, error)
		Insert(context.Context, models.Album) error
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool := dbConnect(ctx)
	defer dbpool.Close()
	log.Println("Connected to database")

	env := &Env{
		Albums: models.AlbumModel{DbPool: dbpool},
	}

	router := gin.Default()
	router.GET("/ping", handlePing)
	router.GET("/albums", env.handleGetAlbums)
	router.GET("/albums/:id", env.handleGetAlbumByID)
	router.POST("/albums", env.handlePostAlbums)

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
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Ok"})
}

// handleGetAlbums responds with the list of all albums as JSON.
func (env *Env) handleGetAlbums(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	albums, err := env.Albums.GetAll(ctx)
	if err != nil {
		c.IndentedJSON(
			http.StatusInternalServerError,
			gin.H{"message": "Internal error", "error": err.Error()},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

// handlePostAlbums adds an album from JSON received in the request body.
func (env *Env) handlePostAlbums(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var a models.Album

	// Call BindJSON to bind the received JSON to a new Album.
	if err := c.BindJSON(&a); err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message": "Invalid request", "error": err.Error()},
		)
		return
	}
	// Ignore ID since it is auto-incremented by the database.
	a.ID = ""

	env.Albums.Insert(ctx, a)
	c.IndentedJSON(http.StatusCreated, a)
}

// handleGetAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func (env *Env) handleGetAlbumByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(
			http.StatusBadRequest,
			gin.H{"message": "Invalid ID number", "error": err.Error()},
		)
		return
	}
	album, err := env.Albums.GetByID(ctx, id)
	if err != nil {
		c.IndentedJSON(
			http.StatusNotFound,
			gin.H{"message": "Album not found", "error": err.Error()},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}
