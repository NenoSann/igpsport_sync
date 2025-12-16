// A go implementation of https://github.com/yihong0618/running_page/blob/master/run_page/igpsport_sync.py
package igpsportsync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const BASE_URL = "https://prod.zh.igpsport.com/service/"
const LOGIN_URL = BASE_URL + "auth/account/login"
const ACTIVITY_URL = BASE_URL + "web-gateway/web-analyze/activity/"
const QUERY_URL = ACTIVITY_URL + "queryMyActivity"
const DOWNLOAD_URL = ACTIVITY_URL + "getDownloadUrl/"

type Config struct {
	Username string
	Password string
	PageSize int
}

type IgpsportSync struct {
	Config      Config
	client      *http.Client
	LoginResult *LoginResult
}

// New creates a new instance of IgpsportSync with the provided configuration.
func (s *IgpsportSync) init(config Config) error {
	s.Config = config
	s.client = &http.Client{
		Timeout: 30 * time.Second,
	}
	err := s.Login()
	if err != nil {
		return err
	}
	return nil
}

// New creates a new instance of IgpsportSync with the provided configuration.
func New(config Config) (*IgpsportSync, error) {
	s := &IgpsportSync{}
	err := s.init(config)

	if err != nil {
		return nil, err
	}
	return s, nil
}

// query activity list
func (s *IgpsportSync) GetActivityList(pageNo int, ext Extension) (resp *ActivityListResponse, err error) {
	if pageNo < 1 {
		return nil, fmt.Errorf("pageNo must be greater than 0")
	}

	// Build query parameters
	query := map[string]string{
		"pageNo":   strconv.Itoa(pageNo),
		"pageSize": strconv.Itoa(s.Config.PageSize),
		"sort":     "1",
		"reqType":  string(ext),
	}

	// Create HTTP request with query parameters
	req, err := http.NewRequest("GET", QUERY_URL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add query parameters
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	// Add authorization header
	s.addAuthHeader(req)

	// Execute request
	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer res.Body.Close()

	// Parse response
	var activityListResp ActivityListResponse
	err = json.NewDecoder(res.Body).Decode(&activityListResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &activityListResp, nil
}

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
func (s *IgpsportSync) DownloadFile(url string, file_name string, ext Extension) ([]byte, error) {
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
func (s *IgpsportSync) DownloadAllActivities(ext Extension, callback DownloadCallback) error {
	if callback == nil {
		return fmt.Errorf("callback function is required")
	}

	page := 1
	for {
		// Get activity list for current page
		resp, err := s.GetActivityList(page, ext)
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
				if !callback(activity) {
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
				if !callback(activity) {
					return nil // Stop downloading
				}
				continue
			}

			// Download the file
			data, err := s.DownloadFile(*downloadURL, strconv.Itoa(row.RideID), ext)
			activity := &DownloadedActivity{
				RideID:    row.RideID,
				Title:     row.Title,
				StartTime: row.StartTime,
				Data:      data,
				Error:     err,
			}

			// Call callback
			if !callback(activity) {
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
// maxConcurrency: maximum number of concurrent downloads (must be > 0)
// callback: called for each downloaded activity, return false to stop
// This method uses worker goroutines to download files in parallel
func (s *IgpsportSync) DownloadAllActivitiesWithConcurrency(ext Extension, maxConcurrency int, callback DownloadCallback) error {
	if callback == nil {
		return fmt.Errorf("callback function is required")
	}

	if maxConcurrency <= 0 {
		return fmt.Errorf("maxConcurrency must be greater than 0")
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
				if !callback(activity) {
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
				if !callback(activity) {
					stopMutex.Lock()
					shouldStop = true
					stopMutex.Unlock()
				}
				continue
			}

			// Download the file
			data, err := s.DownloadFile(*downloadURL, strconv.Itoa(row.RideID), ext)
			activity := &DownloadedActivity{
				RideID:    row.RideID,
				Title:     row.Title,
				StartTime: row.StartTime,
				Data:      data,
				Error:     err,
			}

			// Call callback
			if !callback(activity) {
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
		resp, err := s.GetActivityList(page, ext)
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

// try to login and get access token
func (s *IgpsportSync) Login() error {
	req := map[string]string{
		"appId":    "igpsport-web",
		"username": s.Config.Username,
		"password": s.Config.Password,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	res, err := s.client.Post(LOGIN_URL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("login failed, please check your username and password: %v", err)
	}
	defer res.Body.Close()

	// read response body
	ret, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// unmarshal response body
	retJson := &LoginResponse{}
	err = json.Unmarshal(ret, &retJson)
	if err != nil {
		return fmt.Errorf("error when unmarshalling json: %v", err)
	}

	access_token := retJson.Data.Access_token
	if access_token == "" {
		return fmt.Errorf("login failed, please check your username and password")
	}

	// Save access token for future requests
	s.LoginResult = &retJson.Data

	return nil
}

func (s *IgpsportSync) addAuthHeader(req *http.Request) {
	if s.LoginResult != nil && s.LoginResult.Access_token != "" {
		req.Header.Add("Authorization", "Bearer "+s.LoginResult.Access_token)
	}
}
