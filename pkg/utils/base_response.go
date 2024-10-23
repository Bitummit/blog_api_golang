package utils

import (
	_ "github.com/Bitummit/blog_api_golang/docs"

)

type Response struct {
	Status string `json:"status" expanple:"OK"`
	Error string `json:"error,omitempty"`
}


func OK() Response {
	return Response{
		Status: "OK",
	}
}

func Error(msg string) Response {
	return Response{
		Status: "Error",
		Error: msg,
	}
}