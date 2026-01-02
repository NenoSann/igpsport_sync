// A go implementation of https://github.com/yihong0618/running_page/blob/master/run_page/igpsport_sync.py
package igpsportsync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const BASE_URL = "https://prod.zh.igpsport.com/service/"
const LOGIN_URL = BASE_URL + "auth/account/login"
const ACTIVITY_URL = BASE_URL + "web-gateway/web-analyze/activity/"
const QUERY_URL = ACTIVITY_URL + "queryMyActivity"
const DOWNLOAD_URL = ACTIVITY_URL + "getDownloadUrl/"
const DEFAULT_PAGE_SIZE = 20

type Config struct {
	Username string
	Password string
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
func (s *IgpsportSync) GetActivityList(pageNo int, pageSize int, beginTime string, endTime string) (resp *ActivityListResponse, err error) {
	if pageNo < 1 {
		return nil, fmt.Errorf("pageNo must be greater than 0")
	}

	// Build query parameters
	query := map[string]string{
		"pageNo":    strconv.Itoa(pageNo),
		"pageSize":  strconv.Itoa(pageSize),
		"beginTime": beginTime,
		"endTime":   endTime,
		"sort":      "1",
		"reqType":   "1",
	}

	// Create HTTP request with query parameters
	req, err := http.NewRequest("GET", QUERY_URL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Add query parameters
	q := req.URL.Query()
	for key, value := range query {
		if value == "" {
			continue
		}
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
