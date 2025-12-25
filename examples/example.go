package main

import (
	"fmt"
	"os"
	"sync"

	igpsportsync "github.com/NenoSann/igpsport_sync"
)

func main() {
	// Example 1: Serial Download
	example1SerialDownload()

	// Example 2: Concurrent Download
	example2ConcurrentDownload()

	// Example 3: With Statistics
	example3WithStatistics()
}

// Example 1: Simple serial download
func example1SerialDownload() {
	fmt.Println("=== Example 1: Serial Download ===")

	config := igpsportsync.Config{
		Username: os.Getenv("IGPSPORT_USERNAME"),
		Password: os.Getenv("IGPSPORT_PASSWORD"),
		PageSize: 20,
	}

	client, err := igpsportsync.New(config)

	err = client.DownloadAllActivities(igpsportsync.DownloadOptions{
		Extension: igpsportsync.FIT,
		Callback: func(activity *igpsportsync.DownloadedActivity) bool {
			if activity.Error != nil {
				fmt.Printf("✗ Error downloading activity %d: %v\n", activity.RideID, activity.Error)
				return true // Continue with next activity
			}

			fmt.Printf("✓ Downloaded activity %d: %d bytes\n", activity.RideID, len(activity.Data))
			return true // Continue
		},
	})

	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
	}
	fmt.Println()
}

// Example 2: Concurrent download with worker pool
func example2ConcurrentDownload() {
	fmt.Println("=== Example 2: Concurrent Download (5 workers) ===")

	config := igpsportsync.Config{
		Username: os.Getenv("IGPSPORT_USERNAME"),
		Password: os.Getenv("IGPSPORT_PASSWORD"),
		PageSize: 20,
	}

	client, err := igpsportsync.New(config)
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		return
	}

	err = client.DownloadAllActivitiesWithConcurrency(igpsportsync.DownloadOptions{
		Extension:      igpsportsync.FIT,
		MaxConcurrency: 5,
		Callback: func(activity *igpsportsync.DownloadedActivity) bool {
			if activity.Error != nil {
				fmt.Printf("✗ Error: %v\n", activity.Error)
				return true
			}

			fmt.Printf("✓ Activity %d: %d bytes\n", activity.RideID, len(activity.Data))
			return true
		},
	})

	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
	}
	fmt.Println()
}

// Example 3: Download with statistics
func example3WithStatistics() {
	fmt.Println("=== Example 3: Download with Statistics ===")

	config := igpsportsync.Config{
		Username: os.Getenv("IGPSPORT_USERNAME"),
		Password: os.Getenv("IGPSPORT_PASSWORD"),
		PageSize: 20,
	}

	client, err := igpsportsync.New(config)
	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		return
	}

	stats := struct {
		Total      int
		Successful int
		Failed     int
		TotalBytes int64
		mu         sync.Mutex
	}{}

	err = client.DownloadAllActivitiesWithConcurrency(igpsportsync.DownloadOptions{
		Extension:      igpsportsync.FIT,
		MaxConcurrency: 5,
		Callback: func(activity *igpsportsync.DownloadedActivity) bool {
			stats.mu.Lock()
			defer stats.mu.Unlock()

			stats.Total++

			if activity.Error != nil {
				stats.Failed++
				fmt.Printf("[%d] ✗ Activity %d: %v\n", stats.Total, activity.RideID, activity.Error)
			} else {
				stats.Successful++
				stats.TotalBytes += int64(len(activity.Data))
				fmt.Printf("[%d] ✓ Activity %d: %d bytes\n", stats.Total, activity.RideID, len(activity.Data))
			}

			return true
		},
	})

	if err != nil {
		fmt.Printf("Fatal error: %v\n", err)
		return
	}

	// Print statistics
	fmt.Println("\n=== Statistics ===")
	fmt.Printf("Total: %d\n", stats.Total)
	fmt.Printf("Successful: %d\n", stats.Successful)
	fmt.Printf("Failed: %d\n", stats.Failed)
	fmt.Printf("Total size: %.2f MB\n", float64(stats.TotalBytes)/1024/1024)
}
