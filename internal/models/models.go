package models

type Post struct {
	Id int64 		`json:"id"`
	Title string 	`json:"title"`
	Body string		`json:"body"`
	Author int64 	`json:"author"`
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