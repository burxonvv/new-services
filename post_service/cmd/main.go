package main

import (
	"fmt"
	"net"

	"github.com/burxondv/new-services/post-service/config"
	p "github.com/burxondv/new-services/post-service/genproto/post"
	"github.com/burxondv/new-services/post-service/pkg/db"
	"github.com/burxondv/new-services/post-service/pkg/logger"
	"github.com/burxondv/new-services/post-service/service"
	grpcclient "github.com/burxondv/new-services/post-service/service/grpc_client"

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

	postService := service.NewPostService(connDb, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.PostServicePort)
	if err != nil {
		log.Fatal("failed while listening port: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	p.RegisterPostServiceServer(s, postService)

	log.Info("main: server running",
		logger.String("port", cfg.PostServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed while listening: %v", logger.Error(err))
	}
}
