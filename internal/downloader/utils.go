package downloader

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// getFilenameFromURL extracts the filename from a URL
func getFilenameFromURL(downloadURL string) (string, error) {
	parsedURL, err := url.Parse(downloadURL)
	if err != nil {
		return "", err
	}

	// Get the last part of the path
	filename := filepath.Base(parsedURL.Path)
	if filename == "/" || filename == "." {
		filename = "index.html"
	}

	return filename, nil
}

// createOutputFile creates a new file for writing
func createOutputFile(path string) (*os.File, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}

	// Create the file
	return os.Create(path)
}

// getContentLength gets the content length from the HTTP header
func getContentLength(client *http.Client, url string) (int64, error) {
	resp, err := client.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, ErrHTTPRequestFailed
	}

	return resp.ContentLength, nil
}

// isRangeSupported checks if the server supports range requests
func isRangeSupported(client *http.Client, url string) (bool, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Range", "bytes=0-0")
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusPartialContent, nil
}
