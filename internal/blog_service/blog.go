package blogservice

import (
	"context"
	"fmt"

	"github.com/Bitummit/blog_api_golang/internal/models"
)

type PostQueryFunctions interface {
	NewPost(context.Context, models.Post) (int64, error)
	ListPost(context.Context) ([]models.Post, error)
	GetPost(context.Context, int) (*models.Post, error)
	DeletePost(context.Context, int) error
}


func CreatePostService(storage PostQueryFunctions , post models.Post) (int64, error) {
	
	id, err := storage.NewPost(context.Background(), post)
	if err != nil {
		return id, fmt.Errorf("insertion error %v", err)
	}

	return id, nil
}


func ListPostService(storage PostQueryFunctions) ( []models.Post, error) {
	posts, err := storage.ListPost(context.Background())
	if err != nil {		
		return nil, fmt.Errorf("query error %v", err)
	}

    return posts, nil
}


func GetPostService(storage PostQueryFunctions, id int) (*models.Post, error){

	post, err := storage.GetPost(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("query error %v", err)
	}

	return post, nil
}


func DeletePostService(storage PostQueryFunctions, id int) error{

	if err := storage.DeletePost(context.Background(), id); err != nil {
		return fmt.Errorf("query error %v", err)
	} 
	
	return nil
}
