package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

// Chunk represents a portion of a file to be downloaded
type Chunk struct {
	Start int64
	End   int64
}

// Downloader handles concurrent file downloads
type Downloader struct {
	client      *http.Client
	concurrency int
}

// New creates a new Downloader with the specified concurrency level
func New(concurrency int) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		concurrency: concurrency,
	}
}

// Download downloads a file from the given URL to the specified output path
func (d *Downloader) Download(url, outputPath string) error {
	// Validate URL
	if _, err := http.NewRequest("GET", url, nil); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	// Get content length
	totalSize, err := getContentLength(d.client, url)
	if err != nil {
		return fmt.Errorf("failed to get content length: %w", err)
	}

	// Check if range requests are supported
	rangeSupported, err := isRangeSupported(d.client, url)
	if err != nil {
		return fmt.Errorf("failed to check range support: %w", err)
	}

	// Determine output filename
	if outputPath == "" {
		outputPath, err = getFilenameFromURL(url)
		if err != nil {
			return fmt.Errorf("failed to determine output filename: %w", err)
		}
	}

	// Create output file
	file, err := createOutputFile(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// If range is not supported or file is small, download in a single connection
	if !rangeSupported || totalSize < 1024*1024 { // 1MB threshold
		return d.downloadSingle(url, file)
	}

	// Download using multiple connections
	return d.downloadConcurrent(url, file, totalSize)
}

// downloadSingle downloads a file using a single connection
func (d *Downloader) downloadSingle(url string, file *os.File) error {
	resp, err := d.client.Get(url)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrHTTPRequestFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status code %d", ErrHTTPRequestFailed, resp.StatusCode)
	}

	_, err = io.Copy(file, resp.Body)
	return err
}

// downloadConcurrent downloads a file using multiple concurrent connections
func (d *Downloader) downloadConcurrent(url string, file *os.File, totalSize int64) error {
	chunkSize := totalSize / int64(d.concurrency)
	chunks := make([]Chunk, 0, d.concurrency)

	// Create chunks
	var start int64
	for i := 0; i < d.concurrency; i++ {
		end := start + chunkSize - 1
		if i == d.concurrency-1 {
			end = totalSize - 1 // Last chunk gets the remainder
		}

		chunks = append(chunks, Chunk{Start: start, End: end})
		start = end + 1

		if start >= totalSize {
			break
		}
	}


	var wg sync.WaitGroup
	errChan := make(chan error, len(chunks))

	// Download chunks in parallel
	for _, chunk := range chunks {
		wg.Add(1)
		go func(c Chunk) {
			defer wg.Done()

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				errChan <- fmt.Errorf("failed to create request: %w", err)
				return
			}

			rangeHeader := fmt.Sprintf("bytes=%d-%d", c.Start, c.End)
			req.Header.Set("Range", rangeHeader)

			resp, err := d.client.Do(req)
			if err != nil {
				errChan <- fmt.Errorf("download failed: %w", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusPartialContent {
				errChan <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}

			// Write the chunk to the correct position in the file
			if _, err := file.Seek(c.Start, 0); err != nil {
				errChan <- fmt.Errorf("seek failed: %w", err)
				return
			}

			if _, err := io.Copy(file, resp.Body); err != nil {
				errChan <- fmt.Errorf("failed to write chunk: %w", err)
				return
			}

			errChan <- nil
		}(chunk)
	}

	// Wait for all downloads to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for errors
	for err := range errChan {
		if err != nil {
			return fmt.Errorf("download failed: %w", err)
		}
	}

	return nil
}
