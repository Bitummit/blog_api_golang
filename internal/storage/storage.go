package storage

import "context"


type Post struct {
	Id int64
	Title string
	Body string
	Author string
}

type PostQueryFunctions interface {
	NewPost(context.Context) (int, error)
	ListPost(context.Context) ([]Post, error)
	GetPost(context.Context, int) (*Post, error)
	DeletePost(context.Context, int) error
}