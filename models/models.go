package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Album represents data about a record album.
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Custom AlbumModel type which wraps the connection pool.
type AlbumModel struct {
	DbPool *pgxpool.Pool
}

// Fetch all albums from the database.
func (m AlbumModel) GetAll(c context.Context) ([]Album, error) {
	query := "SELECT id, title, artist, price FROM album"
	rows, err := m.DbPool.Query(c, query)
	if err != nil {
		return nil, fmt.Errorf("Unable to query albums: %w", err)
	}
	defer rows.Close()

	albums, err := pgx.CollectRows(rows, pgx.RowToStructByName[Album])
	if err != nil {
		return nil, fmt.Errorf("Unable to find albums data: %w", err)
	}
	return albums, nil
}

// Fetch an album from the database by its ID.
func (m AlbumModel) GetByID(c context.Context, id int) (*Album, error) {
	query := "SELECT id, title, artist, price FROM album WHERE id = @id"
	args := pgx.NamedArgs{"id": id}
	rows, err := m.DbPool.Query(c, query, args)
	if err != nil {
		return nil, fmt.Errorf("Unable to query albums: %w", err)
	}
	defer rows.Close()

	album, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Album])
	if err != nil {
		return nil, fmt.Errorf("Unable to find album: %w", err)
	}
	return &album, nil
}

// Insert a new album into the database
func (m AlbumModel) Insert(c context.Context, a Album) error {
	query := "INSERT INTO album (title, artist, price) VALUES (@title, @artist, @price)"
	args := pgx.NamedArgs{
		"title":  a.Title,
		"artist": a.Artist,
		"price":  a.Price,
	}
	_, err := m.DbPool.Exec(c, query, args)
	if err != nil {
		return fmt.Errorf("Unable to insert row: %w", err)
	}
	return nil
}
