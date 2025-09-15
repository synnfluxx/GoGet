package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/GoGet/internal/downloader"
)

func main() {
	// Parse command line flags
	concurrent := flag.Int("c", 4, "number of concurrent downloads")
	help := flag.Bool("h", false, "show help")
	output := flag.String("o", "", "output file/directory path")
	flag.Parse()

	if *help || len(flag.Args()) == 0 {
		showHelp()
		return
	}

	// Get URL from command line arguments
	url := flag.Arg(0)

	// Create downloader with specified concurrency
	dl := downloader.New(*concurrent)

	// Start download
	err := dl.Download(url, *output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}
}

func showHelp() {
	helpText := `GoGet - A simple concurrent file downloader

Usage:
  goget [options] <url>

Options:
  -c int
        number of concurrent connections (default 4)
  -o string
        output file/directory path
  -h    show this help message
`
	fmt.Print(helpText)
}
