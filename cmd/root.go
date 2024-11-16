package cmd

import (
	"fmt"
	"stress-test/internal/tester"

	"github.com/spf13/cobra"
)

var (
	url         string
	requests    int
	concurrency int
	debugMode   bool
)

var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "CLI tool to load test a service.",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := tester.NewLoadTester(url, requests, concurrency, debugMode)
		report, err := t.Run()
		if err != nil {
			return err
		}
		fmt.Println(report.Generate())
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVar(&url, "url", "", "Service URL to test.")
	rootCmd.Flags().IntVar(&requests, "requests", 1, "Total number of requests.")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 1, "Concurrent requests.")
	rootCmd.Flags().BoolVar(&debugMode, "debug", false, "Enable debug mode.")

	rootCmd.MarkFlagRequired("url")
}
