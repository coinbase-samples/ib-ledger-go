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

package service

import (
	api "github.com/coinbase-samples/ib-ledger-go/pkg/pbs/ledger/v1"
	"github.com/coinbase-samples/ib-ledger-go/repository"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	api.UnimplementedLedgerServer
	PostgresHandler repository.Repository
	Tracer          trace.Tracer
}

func NewService(pql repository.Repository, tp trace.Tracer) *Service {
	return &Service{PostgresHandler: pql, Tracer: tp}
}
