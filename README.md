# igpsport_sync

A Go library for syncing activities from iGPSport to your application. Download activity data in multiple formats (FIT, GPX, TCX) with support for serial or concurrent downloads.  

This Library is a golang implemet from https://github.com/yihong0618/running_page/blob/master/run_page/igpsport_sync.py

## Features

- **Authentication**: Login with iGPSport account credentials
- **Activity Fetching**: Retrieve activity list with pagination support
- **Multiple Formats**: Support for FIT, GPX, and TCX file formats
- **Flexible Download Options**:
  - Serial download (`DownloadAllActivities`)
  - Concurrent download with customizable worker pool (`DownloadAllActivitiesWithConcurrency`)
- **Callback-based Processing**: Handle each downloaded activity in real-time
- **Error Handling**: Graceful error handling for individual downloads

## Installation

```bash
go get github.com/NenoSann/igpsport_sync
```

## Quick Start

```go
package main

import (
    "fmt"
    igpsportsync "github.com/NenoSann/igpsport_sync"
)

func main() {
    // Create configuration
    config := igpsportsync.Config{
        Username: "your_username",
        Password: "your_password",
        PageSize: 20,
    }

    // Initialize client (automatically logs in)
    client := igpsportsync.New(config)

    // Download activities serially
    client.DownloadAllActivities(igpsportsync.FIT, func(activity *igpsportsync.DownloadedActivity) bool {
        if activity.Error != nil {
            fmt.Printf("Error downloading activity %d: %v\n", activity.RideID, activity.Error)
            return true // Continue with next activity
        }

        fmt.Printf("Downloaded activity %d: %d bytes\n", activity.RideID, len(activity.Data))
        // Process the activity data (e.g., save to file, upload to OSS)
        return true // Continue
    })
}
```

## Concurrent Download Example

For better performance, use concurrent downloads:

```go
// Download with 5 concurrent workers
client.DownloadAllActivitiesWithConcurrency(igpsportsync.FIT, 5, func(activity *igpsportsync.DownloadedActivity) bool {
    if activity.Error != nil {
        fmt.Printf("Error: %v\n", activity.Error)
        return true
    }
    
    // Upload to OSS or save to disk
    saveActivity(activity.Data, activity.RideID)
    return true
})
```

## API Reference

### Types

- `Config`: Configuration for iGPSport client
- `Extension`: File format (FIT, GPX, TCX)
- `DownloadedActivity`: Downloaded activity data with metadata
- `DownloadCallback`: Callback function signature

### Methods

- `New(config Config) *IgpsportSync`: Create a new client
- `GetActivityList(pageNo int, ext Extension) (*ActivityListResponse, error)`: Get activities for a specific page
- `DownloadAllActivities(ext Extension, callback DownloadCallback) error`: Download all activities serially
- `DownloadAllActivitiesWithConcurrency(ext Extension, maxConcurrency int, callback DownloadCallback) error`: Download activities with concurrency

## Documentation

For more detailed information, see:
- [Concurrent Download Guide](./CONCURRENT_DOWNLOAD_GUIDE.md)

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

Based on [running_page](https://github.com/yihong0618/running_page) Python implementation.
