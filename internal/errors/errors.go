package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error type to wrap status.Status.
// This is used as status covers both errors and success codes.
type LedgerError struct {
	Code    codes.Code
	Message string
}

func New(c codes.Code, m string) *LedgerError {
	return &LedgerError{
		Code:    c,
		Message: m,
	}
}

func FromError(err error) *LedgerError {
	s, ok := status.FromError(err)

	if !ok {
		s = status.New(codes.Internal, "failed to convert error")
	}
	return &LedgerError{
		Code:    s.Code(),
		Message: s.Message(),
	}
}

func (l *LedgerError) Error() string {
	return l.Message
}

func (l *LedgerError) ToGrpcError() error {
	return status.Error(l.Code, l.Message)
}
