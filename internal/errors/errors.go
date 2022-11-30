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

// Method used by the `Status.FromError` class to determine the GRPCStatus code
// associated with an error.
// See: https://pkg.go.dev/google.golang.org/grpc/status#FromError
func (e *LedgerError) GRPCStatus() *status.Status {
	if e == nil {
		return nil
	}
	return status.New(e.Code, e.Message)
}
