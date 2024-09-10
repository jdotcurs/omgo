package omgo_test

import (
	"testing"

	"github.com/hectormalot/omgo"
	"github.com/stretchr/testify/require"
)

func TestErrInvalidInput_Error(t *testing.T) {
	err := omgo.ErrInvalidInput{Param: "test", Value: 123}
	require.Equal(t, "invalid input: test = 123", err.Error())
}

func TestErrAPIResponse_Error(t *testing.T) {
	err := omgo.ErrAPIResponse{StatusCode: 400, Message: "Bad Request"}
	require.Equal(t, "API error (status 400): Bad Request", err.Error())
}

func TestErrRateLimit_Error(t *testing.T) {
	err := omgo.ErrRateLimit{Message: "Too many requests"}
	require.Equal(t, "rate limit exceeded: Too many requests", err.Error())
}
