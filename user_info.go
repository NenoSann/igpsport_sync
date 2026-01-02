package igpsportsync

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const USER_INFO_URL = BASE_URL + "mobile/api/User/UserInfo"

// A UserInfo http wrapper function
func (s *IgpsportSync) GetUserInfo() (*UserInfoResponse, error) {
	// Create HTTP request
	req, err := http.NewRequest("GET", USER_INFO_URL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating user info request: %v", err)
	}

	// Add authorization header
	s.addAuthHeader(req)

	// Execute request
	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing user info request: %v", err)
	}
	defer res.Body.Close()

	// Parse response
	var userInfoResp UserInfoResponse
	err = json.NewDecoder(res.Body).Decode(&userInfoResp)
	if err != nil {
		return nil, fmt.Errorf("error decoding user info response: %v", err)
	}

	// Check response code
	if userInfoResp.Code != 0 {
		return nil, fmt.Errorf("user info API error: %s (code: %d)", userInfoResp.Message, userInfoResp.Code)
	}

	return &userInfoResp, nil
}
