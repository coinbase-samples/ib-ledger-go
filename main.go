package main

import (
	"LedgerApp/repository"
	"LedgerApp/service"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	api "LedgerApp/protos/ledger"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

func main() {
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	dba := repository.NewPostgresHandler()
	service := service.NewService(dba)

	server := grpc.NewServer()

	api.RegisterLedgerServer(server, service)
	reflection.Register(server)
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
