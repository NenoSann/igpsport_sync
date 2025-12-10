# Quick Start Guide

Get started with igpsport_sync in 5 minutes.

## Installation

```bash
go get github.com/NenoSann/igpsport_sync
```

## Basic Usage

### 1. Import the Package

```go
import igpsportsync "github.com/NenoSann/igpsport_sync"
```

### 2. Create Configuration

```go
config := igpsportsync.Config{
    Username: "your_igpsport_username",
    Password: "your_igpsport_password",
    PageSize: 20,
}
```

### 3. Create Client (Auto-login)

```go
client := igpsportsync.New(config)
```

### 4. Download Activities

#### Serial Download

```go
err := client.DownloadAllActivities(igpsportsync.FIT, func(activity *igpsportsync.DownloadedActivity) bool {
    if activity.Error != nil {
        fmt.Printf("Error: %v\n", activity.Error)
        return true // Continue
    }
    
    fmt.Printf("Downloaded: %d bytes\n", len(activity.Data))
    // Process activity data here
    return true // Continue
})
```

#### Concurrent Download (Recommended for large datasets)

```go
err := client.DownloadAllActivitiesWithConcurrency(
    igpsportsync.FIT,
    5, // Max 5 concurrent downloads
    func(activity *igpsportsync.DownloadedActivity) bool {
        if activity.Error != nil {
            fmt.Printf("Error: %v\n", activity.Error)
            return true
        }
        
        // Save or process file
        saveFile(activity)
        return true
    },
)
```

## File Formats

- `igpsportsync.FIT` - Binary format, contains detailed metrics
- `igpsportsync.GPX` - GPX format, compatible with most apps
- `igpsportsync.TCX` - TCX format, compatible with Garmin, Training Peaks

## Environment Variables

For security, use environment variables instead of hardcoding credentials:

```bash
export IGPSPORT_USERNAME=your_username
export IGPSPORT_PASSWORD=your_password
```

```go
config := igpsportsync.Config{
    Username: os.Getenv("IGPSPORT_USERNAME"),
    Password: os.Getenv("IGPSPORT_PASSWORD"),
    PageSize: 20,
}
```

## Common Patterns

### Pattern 1: Save to Files

```go
client.DownloadAllActivities(igpsportsync.FIT, func(activity *igpsportsync.DownloadedActivity) bool {
    if activity.Error != nil {
        return true
    }
    
    filename := fmt.Sprintf("activities/%d.fit", activity.RideID)
    os.WriteFile(filename, activity.Data, 0644)
    return true
})
```

### Pattern 2: Upload to Cloud Storage

```go
client.DownloadAllActivitiesWithConcurrency(igpsportsync.FIT, 5, func(activity *igpsportsync.DownloadedActivity) bool {
    if activity.Error != nil {
        return true
    }
    
    // Upload to S3, OSS, or other cloud storage
    uploadToOSS(activity.Data, activity.FileName)
    return true
})
```

### Pattern 3: Process with Progress Tracking

```go
var count, total int
var mu sync.Mutex

client.DownloadAllActivities(igpsportsync.FIT, func(activity *igpsportsync.DownloadedActivity) bool {
    mu.Lock()
    count++
    fmt.Printf("Processing %d/%d\n", count, total)
    mu.Unlock()
    
    process(activity)
    return true
})
```

## Troubleshooting

### Authentication Failed

- Verify username and password are correct
- Check if iGPSport account is active
- Ensure network connectivity

### Timeout Error

- Check network connection
- The default timeout is 30 seconds
- For large files, may need to adjust timeout

### Empty Download URL

- Activity may not have associated files in the requested format
- Try a different format (FIT, GPX, TCX)

## Next Steps

- Read [CONCURRENT_DOWNLOAD_GUIDE.md](./CONCURRENT_DOWNLOAD_GUIDE.md) for advanced usage
- Check [examples/](./examples/) directory for complete examples
- Visit [pkg.go.dev](https://pkg.go.dev/github.com/NenoSann/igpsport_sync) for full API documentation

## Support

For issues, questions, or feature requests, please open an issue on GitHub.
