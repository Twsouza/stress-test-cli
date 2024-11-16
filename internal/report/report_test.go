package report

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReport_Generate(t *testing.T) {
	statusCodes := map[int]int{
		200: 8,
		404: 2,
	}
	r := NewReport(5*time.Second, 10, statusCodes)
	output := r.Generate()

	expected := "\n--- Load Test Report ---\n"
	expected += "Total Time Spent: 5s\n"
	expected += "Total Requests Made: 10\n"
	expected += "Requests with HTTP Status 200: 8\n"
	expected += "Distribution of Other HTTP Status Codes:\n"
	expected += "Status 404: 2\n"

	assert.Equal(t, expected, output)
}
