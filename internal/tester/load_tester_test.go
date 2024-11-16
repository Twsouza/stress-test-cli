package tester

import (
	"io"
	"net/http"
	mocks "stress-test/mocks/tester"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadTester_Run(t *testing.T) {
	mockClient := mocks.NewMockHTTPClient(t)

	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(nil),
	}

	mockClient.On("Get", "http://example.com").Return(mockResponse, nil)

	lt := &LoadTester{
		url:         "http://example.com",
		requests:    10,
		concurrency: 2,
		client:      mockClient,
	}

	report, err := lt.Run()
	assert.NoError(t, err)
	assert.Equal(t, 10, report.TotalRequests)
	assert.Equal(t, 10, report.StatusCodes[200])
}
