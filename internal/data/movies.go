package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"greenlight/internal/validator"
	"time"
)

type MovieModel struct {
	DB *sql.DB
}

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	Runtime   Runtime   `json:"runtime"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
}

func (m MovieModel) Insert(movie *Movie) error {

	querySchema := `
			INSERT INTO movies (title, year, runtime, genres)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version`

	//Notice:movie.Genres is an array so we must use pq.Array to pass as parameter
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, querySchema, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

}

func (m MovieModel) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	querySchema := `
                   SELECT id, created_at, title, year, runtime, genres, version
					FROM movies
					WHERE id = $1`
	var movieResult Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, querySchema, id).Scan(
		&movieResult.ID,
		&movieResult.CreatedAt,
		&movieResult.Title,
		&movieResult.Year,
		&movieResult.Runtime,
		pq.Array(&movieResult.Genres),
		&movieResult.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err

		}
	}
	return &movieResult, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {

	querySchema := `UPDATE movies
					SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
					WHERE id = $5 AND version = $6
					RETURNING version`
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID, movie.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, querySchema, args...).Scan(&movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM movies
              WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsChanged, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsChanged == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func ValidateMovie(v *validator.Validator, input *Movie) {

	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in th efuturee")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genres")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")

	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

}
