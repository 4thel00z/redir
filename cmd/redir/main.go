package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"redir"
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
	urlFlag := flag.String("url", "", "The URL to follow redirections for")
	outputFlag := flag.String("output", "table", "Output format: json or table")
	maxRedirs := flag.Int("max", 10, "Maximum number of redirections to follow")
	flag.Parse()

	if *urlFlag == "" {
		fmt.Println(red + "Error: -url flag is required" + reset)
		flag.Usage()
		os.Exit(1)
	}

	// Execute the redirection follow.
	steps, err := redir.FollowRedirects(*urlFlag, *maxRedirs)
	if err != nil {
		fmt.Println(red+"Error:"+reset, err)
		os.Exit(1)
	}

	// Output according to requested format.
	switch *outputFlag {
	case "json":
		outJSON, err := json.MarshalIndent(steps, "", "  ")
		if err != nil {
			fmt.Println(red+"JSON Marshalling Error:"+reset, err)
			os.Exit(1)
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
	fmt.Printf("\n%sFinished at %s%s\n", rocketEmoji, time.Now().Format(time.RFC1123), reset)
}
