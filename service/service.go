package service

import (
	api "LedgerApp/protos/ledger"
	"LedgerApp/repository"
)

type Service struct {
	api.UnimplementedLedgerServer
	PostgresHandler repository.Repository
}

func NewService(pql repository.Repository) *Service {
	return &Service{PostgresHandler: pql}
}
