package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Config holds the CLI configuration
type Config struct {
	URL         string
	Requests    int
	Concurrency int
}

// Result holds the result of a single HTTP request
type Result struct {
	StatusCode int
	Duration   time.Duration
	Error      error
}

// Report holds the final report data
type Report struct {
	TotalTime     time.Duration
	TotalRequests int
	StatusCodes   map[int]int
	SuccessCount  int
	ErrorCount    int
}

func main() {
	// Parse CLI arguments
	config, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	// Validate configuration
	if err := validateConfig(config); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting load test...\n")
	fmt.Printf("URL: %s\n", config.URL)
	fmt.Printf("Total Requests: %d\n", config.Requests)
	fmt.Printf("Concurrency: %d\n\n", config.Concurrency)

	// Run the load test
	report := runLoadTest(config)

	// Display the report
	displayReport(report)
}

func parseFlags() (*Config, error) {
	var config Config

	flag.StringVar(&config.URL, "url", "", "URL of the service to test (required)")
	flag.IntVar(&config.Requests, "requests", 0, "Total number of requests to make (required)")
	flag.IntVar(&config.Concurrency, "concurrency", 0, "Number of concurrent requests (required)")

	flag.Parse()

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.URL == "" {
		return fmt.Errorf("URL is required. Use --url flag")
	}
	if config.Requests <= 0 {
		return fmt.Errorf("requests must be greater than 0. Use --requests flag")
	}
	if config.Concurrency <= 0 {
		return fmt.Errorf("concurrency must be greater than 0. Use --concurrency flag")
	}
	if config.Concurrency > config.Requests {
		return fmt.Errorf("concurrency cannot be greater than total requests")
	}

	return nil
}

func runLoadTest(config *Config) *Report {
	startTime := time.Now()

	// Create channels for communication
	requestChan := make(chan int, config.Requests)
	resultChan := make(chan Result, config.Requests)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < config.Concurrency; i++ {
		wg.Add(1)
		go worker(&wg, client, config.URL, requestChan, resultChan)
	}

	// Send requests to workers
	go func() {
		for i := 0; i < config.Requests; i++ {
			requestChan <- i
		}
		close(requestChan)
	}()

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results
	report := &Report{
		StatusCodes: make(map[int]int),
	}

	for result := range resultChan {
		report.TotalRequests++

		if result.Error != nil {
			report.ErrorCount++
		} else {
			report.StatusCodes[result.StatusCode]++
			if result.StatusCode == 200 {
				report.SuccessCount++
			}
		}
	}

	report.TotalTime = time.Since(startTime)
	return report
}

func worker(wg *sync.WaitGroup, client *http.Client, url string, requestChan <-chan int, resultChan chan<- Result) {
	defer wg.Done()

	for range requestChan {
		start := time.Now()
		resp, err := client.Get(url)
		duration := time.Since(start)

		result := Result{
			Duration: duration,
			Error:    err,
		}

		if err == nil {
			result.StatusCode = resp.StatusCode
			resp.Body.Close()
		}

		resultChan <- result
	}
}

func displayReport(report *Report) {
	fmt.Println("=== LOAD TEST REPORT ===")
	fmt.Printf("Total execution time: %v\n", report.TotalTime)
	fmt.Printf("Total requests made: %d\n", report.TotalRequests)
	fmt.Printf("Successful requests (HTTP 200): %d\n", report.SuccessCount)

	if report.ErrorCount > 0 {
		fmt.Printf("Failed requests (errors): %d\n", report.ErrorCount)
	}

	fmt.Println("\nHTTP Status Code Distribution:")
	for statusCode, count := range report.StatusCodes {
		percentage := float64(count) / float64(report.TotalRequests) * 100
		fmt.Printf("  %d: %d requests (%.1f%%)\n", statusCode, count, percentage)
	}

	if report.TotalRequests > 0 {
		avgTime := report.TotalTime / time.Duration(report.TotalRequests)
		reqPerSec := float64(report.TotalRequests) / report.TotalTime.Seconds()
		fmt.Printf("\nPerformance Metrics:\n")
		fmt.Printf("  Average request time: %v\n", avgTime)
		fmt.Printf("  Requests per second: %.2f\n", reqPerSec)
	}
}
