package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTinyBench(t *testing.T) {
	tmpDir := t.TempDir()

	testCases := []struct {
		name         string
		filePath     string
		fileContents string
		errorMsg     string
	}{
		{
			name:     "No arguments provided",
			filePath: "",
			errorMsg: "Please provide a path to a valid JS file.",
		},
		{
			name: "No benchmark tags present",
			fileContents: `
				const foo = 'bar';
			`,
			filePath: filepath.Join(tmpDir, "no-benchmark-tags.js"),
			errorMsg: `No benchmarks found. Please define at least one benchmark:`,
		},
		{
			name: "Mismatched benchmark tags",
			fileContents: `
				// tinybench start
				const foo = 'bar';
			`,
			filePath: filepath.Join(tmpDir, "mismatched-benchmark-tags.js"),
			errorMsg: `No benchmarks found. Please define at least one benchmark:`,
		},
		{
			name: "Valid benchmark",
			fileContents: `
				// tinybench start
				const foo = 'bar';
				// tinybench stop
			`,
			filePath: filepath.Join(tmpDir, "valid-benchmark.js"),
			errorMsg: "",
		},
	}

	for _, tc := range testCases {
		if tc.fileContents != "" {
			if err := os.WriteFile(tc.filePath, []byte(tc.fileContents), 0644); err != nil {
				t.Fatal(err)
			}
		}

		t.Run(tc.name, func(t *testing.T) {
			os.Args = []string{"tinybench"}
			if tc.filePath != "" {
				os.Args = append(os.Args, tc.filePath)
			}

			err := run()

			if tc.errorMsg == "" {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected an error with message '%s', but got nil", tc.errorMsg)
				} else if !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("Expected error message '%s', but got '%s'", tc.errorMsg, err.Error())
				}
			}
		})
	}
}
