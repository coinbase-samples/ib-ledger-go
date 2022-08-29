package utils

import api "LedgerApp/protos/ledger"

func GetTransactionTypeFromString(s string) api.TransactionType {
	switch s {
	case "TRANSFER":
		return api.TransactionType_TRANSFER
	default:
		return api.TransactionType_TRANSFER
	}
}

func GetStringFromTransactionType(t api.TransactionType) string {
	switch t {
	case api.TransactionType_TRANSFER:
		return "TRANSFER"
	default:
		return "TRANSFER"
	}
}
