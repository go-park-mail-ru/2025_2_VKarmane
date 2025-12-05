package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	profile := flag.String("profile", "coverage.out", "path to coverage profile")
	threshold := flag.Float64("threshold", 60, "minimum coverage percentage")
	flag.Parse()

	total, covered, err := parseProfile(*profile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse coverage profile: %v\n", err)
		os.Exit(1)
	}

	if total == 0 {
		fmt.Fprintln(os.Stderr, "coverage profile is empty")
		os.Exit(1)
	}

	coverage := (covered / total) * 100
	fmt.Printf("Total coverage: %.2f%% (required %.2f%%)\n", coverage, *threshold)

	if coverage < *threshold {
		os.Exit(1)
	}
}

func parseProfile(path string) (total float64, covered float64, err error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		stmts, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			return 0, 0, err
		}

		count, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return 0, 0, err
		}

		total += stmts
		if count > 0 {
			covered += stmts
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return total, covered, nil
}



