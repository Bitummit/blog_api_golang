package postgresql

import (
	"context"
	"os"
	"time"
	"github.com/Bitummit/blog_api_golang/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func InitDB(ctx context.Context) (*Storage, error) {
	dbURL := os.Getenv("DB_URL")
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)

	defer cancel()

	connPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	if err := connPool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Storage{DB: connPool}, nil
}


func (s *Storage) NewPost(ctx context.Context, post models.Post) (int64, error) {
	query := `
		INSERT INTO post(title, body, author_id) VALUES(@title, @body, @author) RETURNING id;
	`
	args := pgx.NamedArgs{
		"title": post.Title,
		"body": post.Body,
		"author": post.Author,
	}

	var id int64
	err := s.DB.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}


func (s *Storage) ListPost(ctx context.Context) ([]models.Post, error) {
	query := `
		SELECT * FROM post;
	`
	var posts []models.Post

	rows, err := s.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post models.Post

		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.Author)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}


func (s *Storage) GetPost(ctx context.Context, id int) (*models.Post, error) {
	query := `
		SELECT * FROM post WHERE id=@id
	`
	args := pgx.NamedArgs{
		"id": id,
	}
	var post models.Post

	err := s.DB.QueryRow(ctx, query, args).Scan(&post.Id, &post.Title, &post.Body, &post.Author)
	if err != nil {
		return nil, err
	}
	
	return &post, nil
}

func (s *Storage) DeletePost(ctx context.Context, id int) error {
	query := `
		DELETE FROM post WHERE id=@id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := s.DB.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
