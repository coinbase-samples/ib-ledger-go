/**
 * Copyright 2022-present- Present Coinbase Global, Inc.
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
	"github.com/coinbase-samples/ib-ledger-go/internal/service"
	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func gRPCListen(app config.AppConfig, l *log.Entry) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", app.Port))
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %w", err)
	}

	grpcOptions := setupGrpcOptions(app, l)
	s := grpc.NewServer(grpcOptions...)

	setupHealthCheckServer(s)

	api.RegisterLedgerServer(s, &service.Service{App: app})
	reflection.Register(s)

	l.Debugf("gRPC Server starting on port %s\n", app.Port)
	if err := s.Serve(lis); err != nil {
		l.Fatalf("failed to listen for gRPC: %w", err)
	}
}

func setupGrpcOptions(app config.AppConfig, l *log.Entry) []grpc.ServerOption {
	// Logrus entry is used, allowing pre-definition of certain fields by the user.
	// See example setup here https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/logging/logrus/examples_test.go
	opts := []grpcLogrus.Option{
		grpcLogrus.WithDurationField(
			func(duration time.Duration) (key string, value interface{}) {
				return "grpc.time_ns", duration.Nanoseconds()
			},
		),
		grpcLogrus.WithDecider(
			func(fullMethodName string, err error) bool {
				// will not log gRPC calls if it was a call to healthcheck and no error was raised
				if err == nil && fullMethodName == "/grpc.health.v1.Health/Check" {
					return false
				}

				// by default everything will be logged
				return true
			},
		),
	}

	grpcOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				grpcCtxtags.UnaryServerInterceptor(),
				grpcLogrus.UnaryServerInterceptor(l, opts...),
				grpcValidator.UnaryServerInterceptor(),
				grpcRecovery.UnaryServerInterceptor(),
			),
		),
	}

	if !app.IsLocalEnv() {
		// load tls for grpc
		tlsCredentials, err := loadCredentials()
		if err != nil {
			l.Fatalf("Cannot load TLS credentials: %v", err)
		}

		grpcOptions = append(grpcOptions, grpc.Creds(tlsCredentials))
	}

	return grpcOptions
}

func setupHealthCheckServer(s *grpc.Server) {
	//setup health server
	hs := health.NewServer()
	hs.SetServingStatus(
		"grpc.health.v1.Health",
		grpc_health_v1.HealthCheckResponse_SERVING,
	)
	grpc_health_v1.RegisterHealthServer(s, hs)
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
