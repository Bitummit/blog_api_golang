package storage

import "context"


type Post struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	Body string	`json:"body"`
	Author string `json:"author"`
}

type PostQueryFunctions interface {
	NewPost(context.Context, Post) (int64, error)
	ListPost(context.Context) ([]Post, error)
	GetPost(context.Context, int) (*Post, error)
	DeletePost(context.Context, int) error
}