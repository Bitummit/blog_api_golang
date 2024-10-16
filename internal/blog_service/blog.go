package blogservice

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Bitummit/blog_api_golang/internal/models"
	authclient "github.com/Bitummit/go_auth/pkg/auth_client"
	grpcConfig "github.com/Bitummit/go_auth/pkg/config"
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


func LoginService(storage PostQueryFunctions, log *slog.Logger, user models.User) (*string, error){

	client, err := authclient.NewClient(log, grpcConfig.InitConfig())
	if err != nil {
		return nil, err
	}
	token, err := client.Login(user.Username, user.Password)
	if err != nil {
		return nil, err
	}
	
	return &token.Token, nil
}
