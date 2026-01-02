package igpsportsync

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const ACTIVITY_DETAIL_URL = ACTIVITY_URL + "queryActivityDetail/" // + ride_id

func (s *IgpsportSync) GetActivityDetail(ride_id int) (*ActivityDetailResponse, error) {
	// Build URL with ride_id as path parameter (not query parameter)
	url := ACTIVITY_DETAIL_URL + strconv.Itoa(ride_id)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating activity detail request: %v", err)
	}

	// Add authorization header
	s.addAuthHeader(req)

	// Execute request
	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing activity detail request: %v", err)
	}
	defer res.Body.Close()

	// Parse response
	var detailResp ActivityDetailResponse
	err = json.NewDecoder(res.Body).Decode(&detailResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding activity detail response: %v", err)
	}

	// Check response code
	if detailResp.Code != 0 {
		return nil, fmt.Errorf("activity detail API error: %s (code: %d)", detailResp.Message, detailResp.Code)
	}

	return &detailResp, nil
}
