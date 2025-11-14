package main

import (
    "log"
    "net"

    "google.golang.org/grpc"

    // pb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto"
    // "github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/service"
    // "github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository"
)



func main() {
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal("failed to start AuthService %w", err)
	}

	srv := grpc.NewServer()


	srv.Serve(lis)
}