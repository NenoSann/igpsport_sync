package test

import (
	"testing"
)

// TestGetActivityList tests retrieving the activity list
func TestGetActivityList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := CreateTestClient(t)

	t.Log("Retrieving activity list...")
	activityListResp, err := client.GetActivityList(1, 20, "", "")
	if err != nil {
		t.Fatalf("Could not retrieve activity list: %v", err)
	}

	if activityListResp == nil || activityListResp.Data.Rows == nil || len(activityListResp.Data.Rows) == 0 {
		t.Fatalf("No activities found in the response")
	}

	activities := activityListResp.Data.Rows
	t.Logf("Retrieved %d activities from page 1 (total pages: %d)", len(activities), activityListResp.Data.TotalPage)

	// Display activity details
	for i, activity := range activities {
		if i >= 5 { // Show only first 5
			t.Logf("  ... and %d more activities", len(activities)-5)
			break
		}
		t.Logf("  Activity %d: ID=%d, Title=%s, StartTime=%s", i+1, activity.RideID, activity.Title, activity.StartTime)
	}
}

// TestGetActivityDetail tests retrieving detailed information for a specific activity
func TestGetActivityDetail(t *testing.T) {
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
	t.Logf("Testing GetActivityDetail with ride ID: %d", rideID)

	detail, err := client.GetActivityDetail(rideID)
	if err != nil {
		t.Fatalf("Failed to get activity detail: %v", err)
	}

	if detail == nil || detail.Data.RideId == 0 {
		t.Fatalf("Activity detail is empty or invalid")
	}

	t.Logf("Activity detail retrieved successfully:")
	t.Logf("  RideId: %d", detail.Data.RideId)
	t.Logf("  Title: %s", detail.Data.Title)
	t.Logf("  StartTime: %s", detail.Data.StartTime)
	t.Logf("  RideDistance: %d meters", detail.Data.RideDistance)
	t.Logf("  AvgSpeed: %.3f km/h", detail.Data.AvgSpeed)
	t.Logf("  MaxSpeed: %.3f km/h", detail.Data.MaxSpeed)
	t.Logf("  TotalTime: %d seconds", detail.Data.TotalTime)
	t.Logf("  MovingTime: %d seconds", detail.Data.MovingTime)
	t.Logf("  TotalAscent: %d meters", detail.Data.TotalAscent)
	t.Logf("  Device: %s (version: %s)", detail.Data.DeviceInfo.DeviceName, detail.Data.DeviceInfo.SoftwareVersion)
}

// TestTimeRangeFilter tests the time range filtering functionality
func TestTimeRangeFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	client := CreateTestClient(t)

	t.Log("Testing GetActivityList with time range filters...")

	// Test 1: Get all activities (no time filter)
	t.Run("NoTimeFilter", func(t *testing.T) {
		resp, err := client.GetActivityList(1, 20, "", "")
		if err != nil {
			t.Fatalf("Failed to get activities without time filter: %v", err)
		}
		t.Logf("Total activities (no filter): %d", len(resp.Data.Rows))
		if len(resp.Data.Rows) > 0 {
			t.Logf("  First activity: %s (StartTime: %s)", resp.Data.Rows[0].Title, resp.Data.Rows[0].StartTime)
			if len(resp.Data.Rows) > 1 {
				t.Logf("  Last activity: %s (StartTime: %s)",
					resp.Data.Rows[len(resp.Data.Rows)-1].Title,
					resp.Data.Rows[len(resp.Data.Rows)-1].StartTime)
			}
		}
	})

	// Test 2: Filter by specific date range
	t.Run("SpecificDateRange", func(t *testing.T) {
		beginTime := "2024-01-01"
		endTime := "2024-12-31"

		resp, err := client.GetActivityList(1, 20, beginTime, endTime)
		if err != nil {
			t.Logf("Failed to get activities with time filter: %v", err)
			return
		}

		t.Logf("Activities in range %s to %s: %d", beginTime, endTime, len(resp.Data.Rows))

		if len(resp.Data.Rows) > 0 {
			for i, activity := range resp.Data.Rows {
				if i >= 5 { // Show only first 5
					t.Logf("  ... and %d more activities", len(resp.Data.Rows)-5)
					break
				}
				t.Logf("  Activity %d: %s (ID: %d, StartTime: %s)",
					i+1, activity.Title, activity.RideID, activity.StartTime)
			}
		} else {
			t.Logf("No activities found in the specified time range")
		}
	})

	// Test 3: Filter with only begin time
	t.Run("OnlyBeginTime", func(t *testing.T) {
		beginTime := "2024-06-01"

		resp, err := client.GetActivityList(1, 20, beginTime, "")
		if err != nil {
			t.Logf("Failed to get activities with begin time filter: %v", err)
			return
		}

		t.Logf("Activities after %s: %d", beginTime, len(resp.Data.Rows))
		if len(resp.Data.Rows) > 0 {
			t.Logf("  First activity: %s (StartTime: %s)", resp.Data.Rows[0].Title, resp.Data.Rows[0].StartTime)
		}
	})

	// Test 4: Filter with only end time
	t.Run("OnlyEndTime", func(t *testing.T) {
		endTime := "2024-06-30"

		resp, err := client.GetActivityList(1, 20, "", endTime)
		if err != nil {
			t.Logf("Failed to get activities with end time filter: %v", err)
			return
		}

		t.Logf("Activities before %s: %d", endTime, len(resp.Data.Rows))
		if len(resp.Data.Rows) > 0 {
			t.Logf("  Last activity: %s (StartTime: %s)",
				resp.Data.Rows[len(resp.Data.Rows)-1].Title,
				resp.Data.Rows[len(resp.Data.Rows)-1].StartTime)
		}
	})
}

// TestGetActivityDownloadUrl tests getting the download URL for an activity
func TestGetActivityDownloadUrl(t *testing.T) {
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
	t.Logf("Testing GetActivityDownloadUrl with ride ID: %d", rideID)

	downloadURL, err := client.GetActivityDownloadUrl(rideID)
	if err != nil {
		t.Fatalf("Could not retrieve download URL: %v", err)
	}

	if downloadURL == nil || *downloadURL == "" {
		t.Fatalf("Download URL is empty or nil")
	}

	t.Logf("Successfully retrieved download URL: %s", *downloadURL)
}
