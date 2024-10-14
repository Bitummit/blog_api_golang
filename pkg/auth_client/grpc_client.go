package authclient

import (
	"context"
	"log/slog"

	auth_v1 "github.com/Bitummit/blog_api_golang/pkg/auth_v1/proto"
	"github.com/Bitummit/blog_api_golang/pkg/config"
	"github.com/Bitummit/blog_api_golang/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Client auth_v1.AuthClient
	Cfg *config.Config
	Log *slog.Logger
	Conn *grpc.ClientConn
}


func NewClient(log *slog.Logger, cfg *config.Config) (*AuthClient, error) {

	authClient := AuthClient {
		Cfg: cfg,
		Log: log,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	
	conn, err := grpc.NewClient("127.0.0.1:5300", opts...)
	if err != nil {
		return nil, err
	}
	// defer conn.Close()

	client := auth_v1.NewAuthClient(conn)
	authClient.Client = client
	authClient.Conn = conn

	return &authClient, nil
}


func (a *AuthClient) Login(username string, password string) auth_v1.Token {
	request := &auth_v1.BaseUserInformation {
		Username: username,
		Password: password,
	}
	myToken, err := a.Client.Login(context.Background(), request)
	if err != nil {
		a.Log.Error("fail to dial: %v", logger.Err(err))
	}

	return *myToken
}


// client := r;lgmerl
// client.Login()

