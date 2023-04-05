package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	BENCHMARK_DURATION = 10
	ERROR_PREFIX       = " Error: "
	START_DELIMITER    = "// tinybench start"
	STOP_DELIMITER     = "// tinybench stop"
)

type BenchmarkResult struct {
	Min, Max, Median time.Duration
	Iterations       int
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.SetFlags(0)
	fmt.Print("\033[H\033[2J")
	fmt.Println("\n Welcome to tinybench, a tiny Go tool for benchmarking JavaScript code")

	code, err := readJSFile()
	if err != nil {
		return err
	}

	standardCode, benchmarkCode, err := parse(code)
	if err != nil {
		return err
	}

	results, err := executeBenchmarks(benchmarkCode, standardCode)
	if err != nil {
		return err
	}

	displayResults(results)
	return nil
}

func readJSFile() ([]byte, error) {
	if len(os.Args) < 2 {
		return nil, errors.New(ERROR_PREFIX + "Please provide a path to a valid JS file.")
	}

	absPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		return nil, errors.New(ERROR_PREFIX + err.Error())
	}

	file, err := os.ReadFile(absPath)
	if err != nil {
		return nil, errors.New(ERROR_PREFIX + err.Error())
	}

	return file, nil
}

func parse(code []byte) (string, []string, error) {
	scanner := bufio.NewScanner(bytes.NewReader(code))
	var standardCode, currBenchmark strings.Builder
	var benchmarkSegments []string
	inBenchmark := false

	for scanner.Scan() {
		line := scanner.Text()
		switch strings.TrimSpace(line) {
		case START_DELIMITER:
			inBenchmark = true
			currBenchmark.Reset()
		case STOP_DELIMITER:
			inBenchmark = false
			benchmarkSegments = append(benchmarkSegments, currBenchmark.String())
		default:
			if inBenchmark {
				currBenchmark.WriteString(line + "\n")
			} else {
				standardCode.WriteString(line + "\n")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", nil, err
	}
	if len(benchmarkSegments) == 0 {
		return "", nil, errors.New(ERROR_PREFIX + `No benchmarks found. Please define at least one benchmark:
			// tinybench start
			{{ CODE TO BENCHMARK }}
			// tinybench stop`)
	}

	return standardCode.String(), benchmarkSegments, nil
}

func executeBenchmarks(codeSegments []string, code string) ([]BenchmarkResult, error) {
	results := make([]BenchmarkResult, len(codeSegments))
	nodePath, err := exec.LookPath("node")
	if err != nil {
		return nil, errors.New(ERROR_PREFIX + "Node.js not found. Please install Node.js and try again.")
	}

	for i, segment := range codeSegments {
		fmt.Printf("\n Found benchmark %d:\n\n\t```\n\t%s\n\t```\n\n Executing benchmark %d", i+1, strings.ReplaceAll(segment, "\n", "\n\t"), i+1)

		iterations := 0
		executionTimes := make([]time.Duration, 0)
		stopCh := make(chan struct{})
		go func() {
			for {
				select {
				case <-stopCh:
					return
				default:
					cmd := exec.Command(nodePath, "-e", code+segment)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr

					startTime := time.Now()
					if err := cmd.Run(); err != nil {
						log.Fatal(ERROR_PREFIX, err)
					}
					executionTime := time.Since(startTime)
					executionTimes = append(executionTimes, executionTime)
					iterations++
				}
			}
		}()

		for i := 0; i < BENCHMARK_DURATION; i++ {
			fmt.Print(".")
			time.Sleep(1 * time.Second)
		}

		stopCh <- struct{}{}
		fmt.Println(" done!")

		sort.Slice(executionTimes, func(i, j int) bool {
			return executionTimes[i] < executionTimes[j]
		})

		minTime := executionTimes[0]
		maxTime := executionTimes[len(executionTimes)-1]
		medianTime := executionTimes[len(executionTimes)/2]

		results[i] = BenchmarkResult{Min: minTime, Max: maxTime, Median: medianTime, Iterations: iterations}
	}

	return results, nil
}

func displayResults(results []BenchmarkResult) {
	fmt.Println("\n Results")
	fastestBenchmark := 0
	fastestMedian := results[0].Median
	for i, result := range results {
		if result.Median < fastestMedian {
			fastestBenchmark = i
			fastestMedian = result.Median
		}
	}
	headerFormat := " | %-10s | %-10s | %-10s | %-10s | %-10s | %-12s |\n"
	rowFormat := " | %-10d | %-10d | %-10d | %-10d | %-10d | %-12s |\n"

	fmt.Printf(headerFormat, "Benchmark", "Iterations", "Min(ms)", "Max(ms)", "Median(ms)", "Delta")
	fmt.Println(" " + strings.Repeat("-", 81))

	for i, result := range results {
		differenceLabel := "Fastest"
		if i != fastestBenchmark {
			percentageDifference := float64(result.Median-fastestMedian) / float64(fastestMedian) * 100
			differenceLabel = fmt.Sprintf("%.0f%% Slower", percentageDifference)
		}

		fmt.Printf(rowFormat, i+1, result.Iterations,
			result.Min.Milliseconds(),
			result.Max.Milliseconds(),
			result.Median.Milliseconds(),
			differenceLabel)
	}
}
