package igpsportsync

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func (s *IgpsportSync) GetActivityDownloadUrl(ride_id int) (*string, error) {
	// Build URL with ride_id as path parameter (not query parameter)
	url := DOWNLOAD_URL + strconv.Itoa(ride_id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	s.addAuthHeader(req)

	// send http request
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse response
	var downloadUrlResp map[string]any
	err = json.NewDecoder(res.Body).Decode(&downloadUrlResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// the url is the data value
	urlData, ok := downloadUrlResp["data"]
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}
	urlStr, ok := urlData.(string)
	if !ok {
		return nil, fmt.Errorf("download url is not a string")
	}
	return &urlStr, nil
}

// download file from url
func (s *IgpsportSync) DownloadFile(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	s.addAuthHeader(req)

	// send http request
	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// read response body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DownloadCallback is called for each downloaded activity
// Return true to continue, false to stop downloading
type DownloadCallback func(activity *DownloadedActivity) bool

// DownloadAllActivities downloads all activities of a given type with pagination
// It calls the callback function for each downloaded file
// The callback receives a DownloadedActivity with either Data (on success) or Error (on failure)
// Returning false from the callback will stop the download process
func (s *IgpsportSync) DownloadAllActivities(options DownloadOptions) error {
	if options.Callback == nil {
		return fmt.Errorf("callback function is required")
	}

	page := 1
	for {
		// Get activity list for current page
		resp, err := s.GetActivityList(page, DEFAULT_PAGE_SIZE, options.BeginTime, options.EndTime)
		if err != nil {
			return fmt.Errorf("error getting activity list page %d: %v", page, err)
		}

		// Process each activity on this page
		for _, row := range resp.Data.Rows {
			// Get download URL for this activity
			downloadURL, err := s.GetActivityDownloadUrl(row.RideID)
			if err != nil {
				// Call callback with error
				activity := &DownloadedActivity{
					RideID:    row.RideID,
					Title:     row.Title,
					StartTime: row.StartTime,
					Error:     fmt.Errorf("error getting download URL: %v", err),
				}
				if !options.Callback(activity) {
					return nil // Stop downloading
				}
				continue
			}

			if downloadURL == nil || *downloadURL == "" {
				// Call callback with error
				activity := &DownloadedActivity{
					RideID:    row.RideID,
					Title:     row.Title,
					StartTime: row.StartTime,
					Error:     fmt.Errorf("empty download URL"),
				}
				if !options.Callback(activity) {
					return nil // Stop downloading
				}
				continue
			}

			// Download the file
			data, err := s.DownloadFile(*downloadURL)
			activity := &DownloadedActivity{
				RideID:    row.RideID,
				Title:     row.Title,
				StartTime: row.StartTime,
				Data:      data,
				Error:     err,
			}

			// Call callback
			if !options.Callback(activity) {
				return nil // Stop downloading
			}
		}

		// Check if there are more pages
		if page >= resp.Data.TotalPage {
			break
		}
		page++
	}

	return nil
}

// DownloadAllActivitiesWithConcurrency downloads all activities with concurrent control
// This method uses worker goroutines to download files in parallel
// If MaxConcurrency is 0 or not set, defaults to 5
func (s *IgpsportSync) DownloadAllActivitiesWithConcurrency(options DownloadOptions) error {
	if options.Callback == nil {
		return fmt.Errorf("callback function is required")
	}

	// Set default concurrency if not specified
	maxConcurrency := options.MaxConcurrency
	if maxConcurrency <= 0 {
		maxConcurrency = 5
	}

	// Create channels for work distribution and synchronization
	workChan := make(chan ActivityRow) // Channel to distribute activities
	var wg sync.WaitGroup              // WaitGroup to track workers
	shouldStop := false
	stopMutex := &sync.Mutex{}

	// Worker function that downloads activities
	worker := func() {
		defer wg.Done()
		for row := range workChan {
			// Check if we should stop
			stopMutex.Lock()
			if shouldStop {
				stopMutex.Unlock()
				continue
			}
			stopMutex.Unlock()

			// Get download URL for this activity
			downloadURL, err := s.GetActivityDownloadUrl(row.RideID)
			if err != nil {
				// Call callback with error
				activity := &DownloadedActivity{
					RideID:    row.RideID,
					Title:     row.Title,
					StartTime: row.StartTime,
					Error:     fmt.Errorf("error getting download URL: %v", err),
				}
				if !options.Callback(activity) {
					stopMutex.Lock()
					shouldStop = true
					stopMutex.Unlock()
				}
				continue
			}

			if downloadURL == nil || *downloadURL == "" {
				// Call callback with error
				activity := &DownloadedActivity{
					RideID:    row.RideID,
					Title:     row.Title,
					StartTime: row.StartTime,
					Error:     fmt.Errorf("empty download URL"),
				}
				if !options.Callback(activity) {
					stopMutex.Lock()
					shouldStop = true
					stopMutex.Unlock()
				}
				continue
			}

			// Download the file
			data, err := s.DownloadFile(*downloadURL)
			activity := &DownloadedActivity{
				RideID:    row.RideID,
				Title:     row.Title,
				StartTime: row.StartTime,
				Data:      data,
				Error:     err,
			}

			// Call callback
			if !options.Callback(activity) {
				stopMutex.Lock()
				shouldStop = true
				stopMutex.Unlock()
				continue
			}
		}
	}

	// Start worker goroutines
	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go worker()
	}

	// Fetch pages and send work to workers
	page := 1
	for {
		// Check if we should stop
		stopMutex.Lock()
		if shouldStop {
			stopMutex.Unlock()
			break
		}
		stopMutex.Unlock()

		// Get activity list for current page
		resp, err := s.GetActivityList(page, DEFAULT_PAGE_SIZE, options.BeginTime, options.EndTime)
		if err != nil {
			close(workChan)
			wg.Wait()
			return fmt.Errorf("error getting activity list page %d: %v", page, err)
		}

		// Send activities to workers
		for _, row := range resp.Data.Rows {
			stopMutex.Lock()
			if shouldStop {
				stopMutex.Unlock()
				break
			}
			stopMutex.Unlock()

			workChan <- row
		}

		// Check if there are more pages
		if page >= resp.Data.TotalPage {
			break
		}
		page++
	}

	// Close the work channel and wait for workers to finish
	close(workChan)
	wg.Wait()

	return nil
}

// DownloadSingleActivity downloads a single activity by rideId
// It calls the callback function with the downloaded activity data
// The callback receives a DownloadedActivity with either Data (on success) or Error (on failure)
func (s *IgpsportSync) DownloadSingleActivity(rideId int, callback DownloadCallback) error {
	if callback == nil {
		return fmt.Errorf("callback function is required")
	}

	// Get activity detail first to retrieve metadata
	detail, err := s.GetActivityDetail(rideId)
	if err != nil {
		activity := &DownloadedActivity{
			RideID: rideId,
			Error:  fmt.Errorf("error getting activity detail: %v", err),
		}
		callback(activity)
		return err
	}

	downloadUrl := detail.Data.FitUrl

	data, err := s.DownloadFile(downloadUrl)
	activity := &DownloadedActivity{
		RideID:    rideId,
		Title:     detail.Data.Title,
		StartTime: detail.Data.StartTime,
		Data:      data,
		Error:     err,
	}

	// Call callback
	callback(activity)

	return err
}
