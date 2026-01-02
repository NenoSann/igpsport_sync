package test

import (
	"testing"

	igpsportsync "github.com/NenoSann/igpsport_sync"
)

// TestBusinessFlow tests the complete business workflow:
// 1. Login with credentials
// 2. Get user information
// 3. GetActivityList to retrieve activities
// 4. GetActivityDetail for the first activity
// 5. GetActivityDownloadUrl for the first activity
// 6. DownloadFile to download activity data
// 7. DownloadSingleActivity test
// 8. DownloadAllActivitiesWithConcurrency with limited activities
// 9. Test time range filtering
//
// Each step is dependent on the previous step's success.
// If any step fails, the entire test fails and subsequent steps are skipped.
func TestBusinessFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// ===== Step 1: Login with credentials =====
	t.Log("Step 1: Attempting to login...")
	client := CreateTestClient(t)

	if client.LoginResult == nil || client.LoginResult.Access_token == "" {
		t.Fatalf("Step 1 FAILED: Failed to obtain access token after login")
	}

	accessToken := client.LoginResult.Access_token
	t.Logf("Step 1 PASSED: Successfully logged in, access token: %s...", accessToken[:Min(20, len(accessToken))])

	// ===== Step 2: Get user information =====
	t.Log("Step 2: Retrieving user information...")
	userInfo, err := client.GetUserInfo()
	if err != nil {
		t.Fatalf("Step 2 FAILED: Could not retrieve user info: %v", err)
	}

	if userInfo == nil || userInfo.Data.NickName == "" {
		t.Fatalf("Step 2 FAILED: User info is empty or invalid")
	}

	t.Logf("Step 2 PASSED: Retrieved user info for %s (MemberID: %d, Total rides: %d)",
		userInfo.Data.NickName, userInfo.Data.MemberId, userInfo.Data.RideNum)

	// ===== Step 3: GetActivityList to retrieve activities =====
	t.Log("Step 3: Retrieving activity list...")
	activityListResp, err := client.GetActivityList(1, 20, "", "")
	if err != nil {
		t.Fatalf("Step 3 FAILED: Could not retrieve activity list: %v", err)
	}

	if activityListResp == nil || activityListResp.Data.Rows == nil || len(activityListResp.Data.Rows) == 0 {
		t.Fatalf("Step 3 FAILED: No activities found in the response")
	}

	activities := activityListResp.Data.Rows
	t.Logf("Step 3 PASSED: Retrieved %d activities from page 1 (total pages: %d)", len(activities), activityListResp.Data.TotalPage)

	// Display first few activities
	displayCount := Min(3, len(activities))
	for i := 0; i < displayCount; i++ {
		activity := activities[i]
		t.Logf("  Activity %d: ID=%d, Title=%s, StartTime=%s", i+1, activity.RideID, activity.Title, activity.StartTime)
	}

	// Use the first activity for subsequent steps
	firstActivity := activities[0]
	t.Logf("Using first activity for detailed tests: ID=%d, Title=%s", firstActivity.RideID, firstActivity.Title)

	// ===== Step 4: GetActivityDetail for the first activity =====
	t.Logf("Step 4: Retrieving detailed information for activity ID=%d...", firstActivity.RideID)
	detail, err := client.GetActivityDetail(firstActivity.RideID)
	if err != nil {
		t.Fatalf("Step 4 FAILED: Could not retrieve activity detail: %v", err)
	}

	if detail == nil || detail.Data.RideId == 0 {
		t.Fatalf("Step 4 FAILED: Activity detail is empty or invalid")
	}

	t.Logf("Step 4 PASSED: Retrieved activity detail:")
	t.Logf("  Title: %s", detail.Data.Title)
	t.Logf("  Distance: %d meters, AvgSpeed: %.3f km/h, MaxSpeed: %.3f km/h",
		detail.Data.RideDistance, detail.Data.AvgSpeed, detail.Data.MaxSpeed)
	t.Logf("  Device: %s (version: %s)", detail.Data.DeviceInfo.DeviceName, detail.Data.DeviceInfo.SoftwareVersion)

	// ===== Step 5: GetActivityDownloadUrl for the first activity =====
	t.Logf("Step 5: Retrieving download URL for activity ID=%d...", firstActivity.RideID)
	downloadURL, err := client.GetActivityDownloadUrl(firstActivity.RideID)
	if err != nil {
		t.Fatalf("Step 5 FAILED: Could not retrieve download URL: %v", err)
	}

	if downloadURL == nil || *downloadURL == "" {
		t.Fatalf("Step 5 FAILED: Download URL is empty or nil")
	}

	t.Logf("Step 5 PASSED: Successfully retrieved download URL: %s", *downloadURL)

	// ===== Step 6: DownloadFile to download activity data =====
	t.Logf("Step 6: Downloading activity file...")
	fileData, err := client.DownloadFile(*downloadURL)
	if err != nil {
		t.Fatalf("Step 6 FAILED: Could not download file: %v", err)
	}

	if len(fileData) == 0 {
		t.Fatalf("Step 6 FAILED: Downloaded file is empty")
	}

	t.Logf("Step 6 PASSED: Successfully downloaded %d bytes of data", len(fileData))

	// ===== Step 7: Test DownloadSingleActivity =====
	t.Log("Step 7: Testing DownloadSingleActivity...")
	singleDownloaded := false
	err = client.DownloadSingleActivity(firstActivity.RideID, func(activity *igpsportsync.DownloadedActivity) bool {
		if activity.Error != nil {
			t.Logf("Step 7 WARNING: Error in callback: %v", activity.Error)
			return false
		}
		singleDownloaded = true
		t.Logf("  ✓ Downloaded: %s (%d bytes)", activity.Title, len(activity.Data))
		return true
	})

	if err != nil {
		t.Fatalf("Step 7 FAILED: DownloadSingleActivity error: %v", err)
	}

	if !singleDownloaded {
		t.Fatalf("Step 7 FAILED: Activity was not downloaded")
	}

	t.Log("Step 7 PASSED: Successfully used DownloadSingleActivity")

	// ===== Step 8: DownloadAllActivitiesWithConcurrency =====
	t.Log("Step 8: Downloading multiple activities with concurrency...")
	var testLimit = 4 // Limit to first 4 activities for testing
	var count = 0
	option := igpsportsync.DownloadOptions{
		Extension:      igpsportsync.FIT,
		MaxConcurrency: 3,
		Callback: func(activity *igpsportsync.DownloadedActivity) bool {
			if activity.Error != nil {
				t.Logf("  ✗ Error downloading activity %d: %v", activity.RideID, activity.Error)
				return true // Continue with next activity
			}
			count++
			t.Logf("  ✓ Downloaded activity %d: %s (%d bytes)", activity.RideID, activity.Title, len(activity.Data))
			if count >= testLimit {
				t.Logf("  Reached test limit of %d activities", testLimit)
				return false
			}
			return true
		},
	}
	err = client.DownloadAllActivitiesWithConcurrency(option)

	if err != nil {
		t.Fatalf("Step 8 FAILED: Error during concurrent download: %v", err)
	}

	t.Logf("Step 8 PASSED: Successfully downloaded %d activities with concurrency", count)

	// ===== Step 9: Test Time Range Filtering =====
	t.Log("Step 9: Testing time range filtering...")

	// Test with a specific time range
	beginTime := "2024-01-01"
	endTime := "2024-12-31"

	timeFilteredResp, err := client.GetActivityList(1, 20, beginTime, endTime)
	if err != nil {
		t.Logf("Step 9 WARNING: Could not retrieve time-filtered activity list: %v", err)
		// Not fatal - may be no activities in this range
	} else if timeFilteredResp != nil && len(timeFilteredResp.Data.Rows) > 0 {
		t.Logf("Step 9 PASSED: Retrieved %d activities within time range %s to %s",
			len(timeFilteredResp.Data.Rows), beginTime, endTime)

		// Display first few activities in time range
		displayCount := Min(3, len(timeFilteredResp.Data.Rows))
		for i := 0; i < displayCount; i++ {
			activity := timeFilteredResp.Data.Rows[i]
			t.Logf("  Time-filtered Activity %d: ID=%d, Title=%s, StartTime=%s",
				i+1, activity.RideID, activity.Title, activity.StartTime)
		}
	} else {
		t.Logf("Step 9 INFO: No activities found in time range %s to %s", beginTime, endTime)
	}

	// ===== All steps completed successfully =====
	t.Log("")
	t.Log("============================================================")
	t.Log("✓ BUSINESS FLOW TEST COMPLETED SUCCESSFULLY")
	t.Log("============================================================")
	t.Logf("  1. ✓ Successfully authenticated (token: %s...)", accessToken[:Min(20, len(accessToken))])
	t.Logf("  2. ✓ Retrieved user info for %s", userInfo.Data.NickName)
	t.Logf("  3. ✓ Retrieved activity list (%d activities)", len(activities))
	t.Logf("  4. ✓ Retrieved detailed information for activity: %s", detail.Data.Title)
	t.Logf("  5. ✓ Obtained download URL")
	t.Logf("  6. ✓ Downloaded activity file (%d bytes)", len(fileData))
	t.Logf("  7. ✓ Tested DownloadSingleActivity")
	t.Logf("  8. ✓ Downloaded %d activities with concurrency", count)
	t.Logf("  9. ✓ Tested time range filtering functionality")
	t.Log("============================================================")
}
