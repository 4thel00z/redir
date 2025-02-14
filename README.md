# redir

A simple Go library and CLI tool to follow HTTP redirections and output detailed information such as URL, HTTP status code, and request duration. The CLI supports JSON and a pretty-printed table output (with colors and emojis).

## Installation

Install the CLI tool using Go:

```bash
go install github.com/4thel00z/redir@latest
```
Ensure your $GOPATH/bin (or $HOME/go/bin) is in your PATH.
Usage

Once installed, you can run the CLI as follows:
Table Output

```bash
redir -url="https://example.com" -output=table
```

JSON Output

```bash
redir -url="https://example.com" -output=json
```

Limit the number of redirections with the -max flag:

```bash
redir -url="https://example.com" -max=5
```

## STDIN Mode

If the -url flag is omitted, redir will read from STDIN. Each line is scanned for a URI (even if mixed with other text). For example:

```bash
echo "Visit https://example.com for more info" | redir
```

Or reading from a file:

```bash
cat urls.txt | redir -output=json
```

## Usage as a Library

You can also import the redir package into your Go projects to follow HTTP redirections programmatically. Here's a quick example:

```go
package main

import (
	"fmt"
	"log"

	"github.com/4thel00z/redir/redir"
)

func main() {
	// Follow redirections starting at the given URL.
	steps, err := redir.FollowRedirects("https://example.com", 10)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through the redirection steps and print details.
	for _, step := range steps {
		fmt.Printf("URL: %s, Status: %d, Duration: %v\n", step.URL, step.StatusCode, step.Duration)
	}
}
```

## License

This project is licensed under the GPL-3 license.
