/**
 * Copyright 2022 Coinbase Global, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/coinbase-samples/ib-ledger-go/internal/config"
	"github.com/coinbase-samples/ib-ledger-go/internal/dbmanager"
	"github.com/coinbase-samples/ib-ledger-go/internal/repository"
	"github.com/coinbase-samples/ib-ledger-go/internal/service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	log "github.com/sirupsen/logrus"

	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
)

func main() {

	var app config.AppConfig

	config.Setup(&app)
	log.Println("starting app with config", app)

	if app.Env == "" {
		log.Fatalln("no environment name set")
	}

	logLevel, _ := log.ParseLevel(app.LogLevel)
	log.SetLevel(logLevel)

	//setup conn
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", app.Port))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	// See example setup here https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/logging/logrus/examples_test.go
	logrusEntry := log.NewEntry(log.StandardLogger())
	opts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	var server *grpc.Server
	if app.Env == "local" {
		server = grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
				//aw.InterceptorNew(),
				grpc_validator.UnaryServerInterceptor(),
				grpc_recovery.UnaryServerInterceptor(),
			)),
		)
	} else {
		// load tls for grpc
		tlsCredentials, err := loadCredentials()
		if err != nil {
			log.Fatalln("Cannot load TLS credentials: ", err)
		}
		server = grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
				// aw.InterceptorNew(),
				grpc_validator.UnaryServerInterceptor(),
				grpc_recovery.UnaryServerInterceptor(),
			)),
			grpc.Creds(tlsCredentials),
		)

	}

	//setup health server
	healthServer := health.NewServer()
	healthServer.SetServingStatus("grpc.health.v1.Health", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Setup application service
    dbm := dbmanager.NewPostgresDBManager(app)
	rep := repository.NewPostgresHandler(dbm)
	service := service.NewService(rep)

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
