package test

import (
	"testing"

	igpsportsync "github.com/NenoSann/igpsport_sync"
)

// TestDownloadFile tests downloading a single activity file
func TestDownloadFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := CreateTestClient(t)

	// First get activity list to get a valid ride ID
	activityListResp, err := client.GetActivityList(1, 1, "", "")
	if err != nil {
		t.Fatalf("Could not retrieve activity list: %v", err)
	}

	if len(activityListResp.Data.Rows) == 0 {
		t.Skip("No activities available for testing")
	}

	rideID := activityListResp.Data.Rows[0].RideID
	t.Logf("Testing download with ride ID: %d", rideID)

	// Get download URL
	downloadURL, err := client.GetActivityDownloadUrl(rideID)
	if err != nil {
		t.Fatalf("Could not retrieve download URL: %v", err)
	}

	if downloadURL == nil || *downloadURL == "" {
		t.Fatalf("Download URL is empty or nil")
	}

	// Download the file
	fileData, err := client.DownloadFile(*downloadURL)
	if err != nil {
		t.Fatalf("Could not download file: %v", err)
	}

	if len(fileData) == 0 {
		t.Fatalf("Downloaded file is empty")
	}

	t.Logf("Successfully downloaded %d bytes of data", len(fileData))
}

// `TestDownloadSingleActivity tests the DownloadSingleActivity function
func TestDownloadSingleActivity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := CreateTestClient(t)

	// First get activity list to get a valid ride ID
	activityListResp, err := client.GetActivityList(1, 1, "", "")
	if err != nil {
		t.Fatalf("Could not retrieve activity list: %v", err)
	}

	if len(activityListResp.Data.Rows) == 0 {
		t.Skip("No activities available for testing")
	}

	rideID := activityListResp.Data.Rows[0].RideID
	t.Logf("Testing DownloadSingleActivity with ride ID: %d", rideID)

	downloaded := false
	err = client.DownloadSingleActivity(rideID, func(activity *igpsportsync.DownloadedActivity) bool {
		if activity.Error != nil {
			t.Fatalf("Error downloading activity: %v", activity.Error)
			return false
		}

		t.Logf("Successfully downloaded activity:")
		t.Logf("  RideID: %d", activity.RideID)
		t.Logf("  Title: %s", activity.Title)
		t.Logf("  StartTime: %s", activity.StartTime)
		t.Logf("  Size: %d bytes", len(activity.Data))

		downloaded = true
		return true
	})

	if err != nil {
		t.Fatalf("DownloadSingleActivity failed: %v", err)
	}

	if !downloaded {
		t.Fatalf("Callback was not called")
	}
}

// TestDownloadAllActivities tests the sequential download of multiple activities
func TestDownloadAllActivities(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := CreateTestClient(t)

	t.Log("Testing DownloadAllActivities (sequential)...")

	downloadCount := 0
	maxDownloads := 3 // Limit to 3 activities for testing

	options := igpsportsync.DownloadOptions{
		Extension: igpsportsync.FIT,
		Callback: func(activity *igpsportsync.DownloadedActivity) bool {
			if activity.Error != nil {
				t.Logf("✗ Error downloading activity %d: %v", activity.RideID, activity.Error)
				return true // Continue with next activity
			}

			downloadCount++
			t.Logf("✓ Downloaded activity %d: %s (%d bytes)", activity.RideID, activity.Title, len(activity.Data))

			if downloadCount >= maxDownloads {
				t.Logf("Reached test limit of %d activities", maxDownloads)
				return false // Stop downloading
			}
			return true
		},
	}

	err := client.DownloadAllActivities(options)
	if err != nil {
		t.Fatalf("DownloadAllActivities failed: %v", err)
	}

	t.Logf("Successfully downloaded %d activities", downloadCount)
}

// TestDownloadAllActivitiesWithConcurrency tests the concurrent download of multiple activities
func TestDownloadAllActivitiesWithConcurrency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := CreateTestClient(t)

	t.Log("Testing DownloadAllActivitiesWithConcurrency...")

	downloadCount := 0
	maxDownloads := 5 // Limit to 5 activities for testing

	options := igpsportsync.DownloadOptions{
		Extension:      igpsportsync.FIT,
		MaxConcurrency: 3, // Use 3 concurrent workers
		Callback: func(activity *igpsportsync.DownloadedActivity) bool {
			if activity.Error != nil {
				t.Logf("✗ Error downloading activity %d: %v", activity.RideID, activity.Error)
				return true // Continue with next activity
			}

			downloadCount++
			t.Logf("✓ Downloaded activity %d: %s (%d bytes)", activity.RideID, activity.Title, len(activity.Data))

			if downloadCount >= maxDownloads {
				t.Logf("Reached test limit of %d activities", maxDownloads)
				return false // Stop downloading
			}
			return true
		},
	}

	err := client.DownloadAllActivitiesWithConcurrency(options)
	if err != nil {
		t.Fatalf("DownloadAllActivitiesWithConcurrency failed: %v", err)
	}

	t.Logf("Successfully downloaded %d activities with concurrency", downloadCount)
}

// TestDownloadWithTimeFilter tests downloading with time range filters
func TestDownloadWithTimeFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := CreateTestClient(t)

	t.Log("Testing download with time range filter...")

	beginTime := "2024-01-01"
	endTime := "2024-12-31"
	downloadCount := 0
	maxDownloads := 2

	options := igpsportsync.DownloadOptions{
		Extension: igpsportsync.FIT,
		BeginTime: beginTime,
		EndTime:   endTime,
		Callback: func(activity *igpsportsync.DownloadedActivity) bool {
			if activity.Error != nil {
				t.Logf("  Error downloading activity %d: %v", activity.RideID, activity.Error)
				return true
			}
			downloadCount++
			t.Logf("  ✓ Downloaded activity %d: %s (%d bytes)",
				activity.RideID, activity.Title, len(activity.Data))

			// Stop after maxDownloads
			return downloadCount < maxDownloads
		},
	}

	err := client.DownloadAllActivities(options)
	if err != nil {
		t.Logf("Download with time filter completed: %v", err)
	}

	t.Logf("Successfully downloaded %d activities within time range %s to %s", downloadCount, beginTime, endTime)
}
