package models

import (
	_ "github.com/Bitummit/blog_api_golang/docs"

)

//
type Post struct {
	// The UUID of a post
	Id int64 		`json:"id" example:"345"`
	// Post title
	Title string 	`json:"title" example:"first post"`
	// Post text
	Body string		`json:"body" example:"this is my first psot and ..."`
	// Post author
	Author int64 	`json:"author" example:"234"`
}

type Author struct {
	Id int64 			`json:"id"`
	UserId int64 		`json:"user_id"`
	FirstName string	`json:"first_name"`
	LastName int64 		`json:"last_name"`
	Age int64 			`json:"age"`
}

type User struct {
	Id int64 			`json:"id"`
	Username string 	`json:"username"`
	Password string		`json:"password"`
}