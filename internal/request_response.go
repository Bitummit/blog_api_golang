package internal

import (
	"github.com/Bitummit/blog_api_golang/pkg/utils"
)

type CreatePostRequest struct{
	Title string `json:"title" validate:"required"`
	Body string `json:"body" validate:"required"`
	Author int64 `json:"author" validate:"required"`
}

type ListPostResponse struct{
	Response utils.Response `json:"response"`
	Posts []Post
}

type GetPostResponse struct {
	Response utils.Response `json:"response"`
	Post *Post `json:"post"`
}
