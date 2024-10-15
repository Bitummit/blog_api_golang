package api

import (
	"github.com/Bitummit/blog_api_golang/internal/models"
	"github.com/Bitummit/blog_api_golang/pkg/utils"
)

type CreatePostRequest struct{
	Title string 	`json:"title" validate:"required"`
	Body string 	`json:"body" validate:"required"`
	Author int64 	`json:"author" validate:"required"`
}

type LoginRequest struct{
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct{
	Response utils.Response `json:"response"`
	Token string 			`json:"token"`
}

type ListPostResponse struct{
	Response utils.Response `json:"response"`
	Posts []models.Post		`json:"posts"`
}

type GetPostResponse struct {
	Response utils.Response `json:"response"`
	Post *models.Post 		`json:"post"`
}
