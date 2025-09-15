package downloader

import "fmt"

// ErrInvalidURL is returned when the provided URL is invalid
var ErrInvalidURL = fmt.Errorf("invalid URL")

// ErrHTTPRequestFailed is returned when the HTTP request fails
var ErrHTTPRequestFailed = fmt.Errorf("HTTP request failed")

// ErrInvalidRangeSupport is returned when the server doesn't support range requests
var ErrInvalidRangeSupport = fmt.Errorf("server does not support range requests")

// ErrDownloadFailed is returned when the download fails
var ErrDownloadFailed = fmt.Errorf("download failed")
