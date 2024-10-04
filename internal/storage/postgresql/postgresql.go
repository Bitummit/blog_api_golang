package postgresql

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB *pgxpool.Pool
}

func InitDB(ctx context.Context) (*Storage, error) {
	dbURL := os.Getenv("DB_URL")
	connPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	return &Storage{DB: connPool}, nil
}


// get all posts
// get post
// create post (auth required)
// delete post (auth required)

// get author(future user, sign in)
// create author (sign up)
//	


// func (s *Storage) 