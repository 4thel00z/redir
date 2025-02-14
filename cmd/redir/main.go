package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/4thel00z/redir"
)

// ANSI escape codes for colors.
const (
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	reset  = "\033[0m"
)

// emoji constants
const (
	rocketEmoji = "🚀"
	arrowEmoji  = "➡️"
	checkEmoji  = "✅"
)

func main() {
	// CLI flags
	urlFlag := flag.String("url", "", "The URL to follow redirections for. If empty, URIs will be read from STDIN (one per line, extra text allowed).")
	outputFlag := flag.String("output", "table", "Output format: json or table")
	maxRedirs := flag.Int("max", 10, "Maximum number of redirections to follow")
	flag.Parse()

	// If -url flag is provided, process that URI.
	if *urlFlag != "" {
		processURI(*urlFlag, *maxRedirs, *outputFlag)
		return
	}

	// Otherwise, read from STDIN.
	scanner := bufio.NewScanner(os.Stdin)
	// Regex to match URIs (basic http/https pattern).
	re := regexp.MustCompile(`https?://[^\s]+`)
	for scanner.Scan() {
		line := scanner.Text()
		uri := re.FindString(line)
		if uri == "" {
			continue // Skip lines with no URI.
		}
		processURI(uri, *maxRedirs, *outputFlag)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, red+"Error reading from STDIN:"+reset, err)
		os.Exit(1)
	}
}

func processURI(uri string, maxRedirs int, outputFormat string) {
	steps, err := redir.FollowRedirects(uri, maxRedirs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%sError processing %s: %v%s\n", red, uri, err, reset)
		return
	}

	switch outputFormat {
	case "json":
		outJSON, err := json.MarshalIndent(steps, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sJSON Marshalling Error for %s: %v%s\n", red, uri, err, reset)
			return
		}
		fmt.Println(string(outJSON))
	case "table":
		printTable(steps)
	default:
		fmt.Println(red + "Unknown output format. Use 'json' or 'table'." + reset)
		os.Exit(1)
	}
}

func printTable(steps []redir.Redirection) {
	// Print header.
	fmt.Printf("%s %-3s %-50s %-12s %-15s %s\n", yellow, "#", "URL", "Status", "Duration", reset)
	// Iterate over steps.
	for i, step := range steps {
		emoji := arrowEmoji
		if i == len(steps)-1 {
			emoji = checkEmoji
		}
		ms := strconv.FormatFloat(float64(step.Duration.Milliseconds()), 'f', 0, 64) + "ms"
		statusColor := yellow
		if step.StatusCode < 300 || step.StatusCode >= 400 {
			statusColor = green
		}
		fmt.Printf("%s %-3d %s %-50s %s %-12d %s %-15s %s\n",
			green, i+1, reset,
			step.URL,
			emoji,
			step.StatusCode,
			statusColor, ms, reset)
	}
	fmt.Printf("\n%s Finished at %s%s\n", rocketEmoji, time.Now().Format(time.RFC1123), reset)
}
