package main

import (
	"LedgerApp/repository"
	"LedgerApp/service"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	api "LedgerApp/protos/ledger"
)

var (
	port = flag.Int("port", 8443, "The server port")
)

func main() {
	lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	dba := repository.NewPostgresHandler()
	service := service.NewService(dba)

	tlsCredentials, err := loadCredentials()
	if err != nil {
		log.Fatalln("Cannot load TLS credentials: ", err)
	}

	server := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)

	api.RegisterLedgerServer(server, service)
	reflection.Register(server)
	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func loadCredentials() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		return nil, err
	}

	return credentials.NewTLS(
		&tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
		},
	), nil
}
