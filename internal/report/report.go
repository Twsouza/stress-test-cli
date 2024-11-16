package report

import (
	"bytes"
	"text/template"
	"time"
)

type Report struct {
	ElapsedTime   time.Duration
	TotalRequests int
	StatusCodes   map[int]int
}

func NewReport(elapsedTime time.Duration, totalRequests int, statusCodes map[int]int) *Report {
	return &Report{
		ElapsedTime:   elapsedTime,
		TotalRequests: totalRequests,
		StatusCodes:   statusCodes,
	}
}

func (r *Report) Generate() string {
	reportTemplate := `
--- Load Test Report ---
Total Time Spent: {{ .ElapsedTime }}
Total Requests Made: {{ .TotalRequests }}
Requests with HTTP Status 200: {{ index .StatusCodes 200 }}
{{- if .OtherStatusCodes }}
Distribution of Other HTTP Status Codes:
{{- range $code, $count := .OtherStatusCodes }}
Status {{ $code }}: {{ $count }}
{{- end }}
{{- end }}
{{- if gt .FailedRequests 0 }}
Failed Requests (No HTTP Response): {{ .FailedRequests }}
{{- end }}
`

	data := struct {
		ElapsedTime      time.Duration
		TotalRequests    int
		StatusCodes      map[int]int
		OtherStatusCodes map[int]int
		FailedRequests   int
	}{
		ElapsedTime:   r.ElapsedTime,
		TotalRequests: r.TotalRequests,
		StatusCodes:   r.StatusCodes,
		OtherStatusCodes: func() map[int]int {
			other := make(map[int]int)
			for code, count := range r.StatusCodes {
				if code != 200 && code != 0 {
					other[code] = count
				}
			}
			return other
		}(),
		FailedRequests: r.StatusCodes[0],
	}

	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return "Error generating report"
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "Error generating report"
	}

	return buf.String()
}
