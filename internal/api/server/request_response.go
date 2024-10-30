package server

import (
	"github.com/Bitummit/blog_api_golang/internal/models"
	"github.com/Bitummit/blog_api_golang/pkg/utils"
)

type (
	CreatePostRequest struct{
		Title string 	`json:"title" validate:"required"`
		Body string 	`json:"body" validate:"required"`
		Author int64 	`json:"author" validate:"required"`
	}

	LoginRequest struct{
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RegisterRequest struct{
		Username string `json:"username" validate:"required"`
		Email string	`json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	RegisterResponse struct{
		Response utils.Response `json:"response"`
		Token string 			`json:"token"`
	}

	LoginResponse struct{
		Response utils.Response `json:"response"`
		Token string 			`json:"token"`
	}

	ListPostResponse struct{
		Response utils.Response `json:"response"`
		Posts []models.Post		`json:"posts"`
	}

	GetPostResponse struct {
		Response utils.Response `json:"response"`
		Post *models.Post 		`json:"post"`
	}
)