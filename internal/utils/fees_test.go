package utils

import (
	"testing"

	ledgererr "github.com/coinbase-samples/ib-ledger-go/internal/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
)

func TestBadCurrency(t *testing.T) {
	expectedErr := ledgererr.New(codes.InvalidArgument, "invalid currency for fee accounts: ETH")
	_, _, err := GetFeeAccounts("ETH")
	if assert.Error(t, err) {
		assert.Equal(t, expectedErr, err)
	} else {
		assert.FailNow(t, "expcted error from function execution")
	}
}
