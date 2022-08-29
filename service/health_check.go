package service

import (
	api "LedgerApp/protos/ledger"
	"context"
)

func (s *Service) HealthCheck(ctx context.Context, req *api.HealthCheckRequest) (*api.HealthCheckResponse, error) {
	return &api.HealthCheckResponse{Status: api.HealthCheckResponse_SERVING}, nil
}
