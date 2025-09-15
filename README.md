# GoGet - Concurrent File Downloader

GoGet is a simple, concurrent file downloader written in Go. It allows you to download files from the internet with multiple concurrent connections for faster downloads.

## Features

- Concurrent downloads using multiple connections
- Automatic fallback to single connection if server doesn't support range requests
- Progress reporting
- Customizable output path
- Clean and modular codebase

## Installation

```bash
go install GoGet
```

## Usage

```bash
# Basic usage
goget https://example.com/file.zip

# Specify output file
goget -o output.zip https://example.com/file.zip

# Specify number of concurrent connections (default: 4)
goget -c 8 https://example.com/largefile.iso
```

## Options

```
  -c int
        number of concurrent connections (default 4)
  -o string
        output file/directory path
  -h    show help message
```

## Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/GoGet.git
   cd GoGet
   ```

2. Build the binary:
   ```bash
   go build -o goget
   ```

3. Install it:
   ```bash
   sudo mv goget /usr/local/bin/
   ```

## License

MIT
