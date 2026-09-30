package healthcheck

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testItemOK struct{}

func (t *testItemOK) IsReady() error {
	return nil
}

type testItemFail struct{}

func (t *testItemFail) IsReady() error {
	return errors.New("test item fails")
}

func TestHealthCheckRun(t *testing.T) {
	tests := []struct {
		name           string
		items          []Checking
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "one item ok",
			items:          []Checking{&testItemOK{}},
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
		{
			name:           "one item fails",
			items:          []Checking{&testItemFail{}},
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "test item fails",
		},
		{
			name:           "mixed items fails",
			items:          []Checking{&testItemOK{}, &testItemFail{}},
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "test item fails",
		},
	}

	for idx, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			port := strconv.Itoa(5555 + idx)

			statusCode, body := runHealthCheck(t, test.items, port)

			require.Equal(t, test.expectedStatus, statusCode)
			require.Contains(t, body, test.expectedBody)
		})
	}
}

// runHealthCheck makes request to healthcheck and returns response code and body.
func runHealthCheck(t *testing.T, items []Checking, port string) (int, string) {
	t.Helper()

	health := New(items, WithPort(port))
	errChan := make(chan error, 1)
	// start server
	runCtx, runCancel := context.WithCancel(context.Background())
	t.Cleanup(runCancel)
	go health.Run(runCtx, errChan)
	<-health.Started()

	// wait for true healthcheck exiting
	t.Cleanup(func() {
		runCancel()
		<-health.Exited()
	})

	// create request
	reqCtx, reqCancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(reqCancel)
	addr := "http://127.0.0.1:" + port + "/health"
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, addr, nil)
	require.NoError(t, err)

	// do request
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	// read response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp.StatusCode, string(body)
}
