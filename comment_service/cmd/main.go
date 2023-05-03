package main

import (
	"fmt"
	"net"

	"github.com/burxondv/new-services/comment-service/config"
	c "github.com/burxondv/new-services/comment-service/genproto/comment"
	"github.com/burxondv/new-services/comment-service/pkg/db"
	"github.com/burxondv/new-services/comment-service/pkg/logger"
	"github.com/burxondv/new-services/comment-service/service"
	grpcclient "github.com/burxondv/new-services/comment-service/service/grpc_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "golang")
	defer logger.Cleanup(log)

	connDb, err := db.ConnectToDB(cfg)
	if err != nil {
		fmt.Println("failed connect database", err)
	}

	grpcClient, err := grpcclient.New(cfg)
	if err != nil {
		fmt.Println("failed while grpc client", err.Error())
	}

	commentService := service.NewCommentService(connDb, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.CommentServicePort)
	if err != nil {
		log.Fatal("failed while listening: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	c.RegisterCommentServiceServer(s, commentService)

	log.Info("main: server running",
		logger.String("port", cfg.CommentServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed while listening: %v", logger.Error(err))
	}
}
