package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"reminderai/model"
)

// BookRepository handles database operations for books
type BookRepository struct {
	pool *pgxpool.Pool
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{pool}
}

func (r *BookRepository) Create(book *model.Book) error {
	ctx := context.Background()
	return r.pool.QueryRow(ctx,
		"INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id, created_at, updated_at",
		book.Title, book.Author).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
}

func (r *BookRepository) GetAll() ([]model.Book, error) {
	ctx := context.Background()
	rows, err := r.pool.Query(ctx, "SELECT id, title, author, created_at, updated_at FROM books ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []model.Book
	for rows.Next() {
		var book model.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.CreatedAt, &book.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *BookRepository) Update(book *model.Book) error {
	ctx := context.Background()
	return r.pool.QueryRow(ctx,
		"UPDATE books SET title = $1, author = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING created_at, updated_at",
		book.Title, book.Author, book.ID).Scan(&book.CreatedAt, &book.UpdatedAt)
}
