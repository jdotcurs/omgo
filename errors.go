package omgo

import "fmt"

// ErrInvalidInput represents an error due to invalid input parameters
type ErrInvalidInput struct {
	Param string
	Value interface{}
}

func (e ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input: %s = %v", e.Param, e.Value)
}

// ErrAPIResponse represents an error returned by the Open-Meteo API
type ErrAPIResponse struct {
	StatusCode int
	Message    string
}

func (e ErrAPIResponse) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// ErrRateLimit represents an error due to exceeding the rate limit
type ErrRateLimit struct {
	Message string
}

func (e ErrRateLimit) Error() string {
	return fmt.Sprintf("rate limit exceeded: %s", e.Message)
}
