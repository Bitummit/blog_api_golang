package blogservice

import (
	"context"
	"log/slog"

	"github.com/Bitummit/blog_api_golang/internal/api"
	"github.com/Bitummit/blog_api_golang/internal/models"
	"github.com/Bitummit/blog_api_golang/pkg/logger"
	"github.com/Bitummit/blog_api_golang/pkg/utils"
)

// return http error code
func ListPostService(log *slog.Logger, storage api.PostQueryFunctions) ( []models.Post, utils.Response) {
	
	posts, err := storage.ListPost(context.Background())
	if err != nil {
		log.Error("query error on fetching posts", logger.Err(err))
		
		return nil, utils.Error("server error")
	}

    return posts, utils.OK()
}