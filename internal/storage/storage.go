package storage

import "context"


type Post struct {
	Id int64
	Title string
	Body string
	Author string
}

type PostQueryFunctions interface {
	NewPost(context.Context, Post) (int64, error)
	ListPost(context.Context) ([]Post, error)
	GetPost(context.Context, int) (*Post, error)
	DeletePost(context.Context, int) error
}