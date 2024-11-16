package tester

import (
	"errors"
	"log"
	"net/http"
	"stress-test/internal/report"
	"sync"
	"time"
)

type LoadTester struct {
	url         string
	requests    int
	concurrency int
	debugMode   bool
	client      HTTPClient
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

func NewLoadTester(url string, requests int, concurrency int, debug bool) *LoadTester {
	return &LoadTester{
		url:         url,
		requests:    requests,
		concurrency: concurrency,
		debugMode:   debug,
		client:      &http.Client{},
	}
}

func (lt *LoadTester) Debug(message string, args ...any) {
	if lt.debugMode {
		log.Printf(message, args...)
	}
}

func (lt *LoadTester) Run() (*report.Report, error) {
	if lt.requests <= 0 || lt.concurrency <= 0 {
		return nil, errors.New("requests and concurrency must be greater than 0")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	requestCh := make(chan int)
	statusCodes := make(map[int]int)
	startTime := time.Now()

	lt.Debug("Starting load test")

	for i := 0; i < lt.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requestCh {
				lt.Debug("Sending request")

				resp, err := lt.client.Get(lt.url)
				if err != nil {
					lt.Debug("Error sending request to %s: %s", lt.url, err)

					mu.Lock()
					statusCodes[0]++
					mu.Unlock()

					continue
				}

				lt.Debug("Received response with status code %d", resp.StatusCode)

				mu.Lock()
				statusCodes[resp.StatusCode]++
				mu.Unlock()

				resp.Body.Close()
			}
		}()
	}

	go func() {
		for i := 0; i < lt.requests; i++ {
			requestCh <- i
		}
		close(requestCh)
	}()

	lt.Debug("Waiting for all requests to finish")
	wg.Wait()

	elapsedTime := time.Since(startTime)
	lt.Debug("Load test finished in %s", elapsedTime)

	return report.NewReport(elapsedTime, lt.requests, statusCodes), nil
}
