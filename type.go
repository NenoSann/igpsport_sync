package igpsportsync

type Extension string

const (
	FIT Extension = "0"
	GPX Extension = "1"
	TCX Extension = "2"
)

type ActivityRow struct {
	RideID       int     `json:"rideId"`
	Title        string  `json:"title"`
	RideDistance float64 `json:"rideDistance"`
	StartTime    string  `json:"startTime"`
}

type ActivityListData struct {
	Rows      []ActivityRow `json:"rows"`
	TotalPage int           `json:"totalPage"`
	PageNo    int           `json:"pageNo"`
	PageSize  int           `json:"pageSize"`
	TotalRows int           `json:"totalRows"`
}

type ActivityListResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    ActivityListData `json:"data"`
}

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    LoginResult
}

type LoginResult struct {
	Token_type    string `json:"token_type"`
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Expires_in    int    `json:"expires_in"`
	Scope         string `json:"scope"`
	BoundPhone    bool   `json:"boundPhone"`
}

type DownloadedActivity struct {
	RideID    int
	Title     string
	StartTime string
	Data      []byte
	Error     error
}

// DownloadOptions contains configuration for downloading activities
type DownloadOptions struct {
	// Extension specifies the file format (FIT, GPX, TCX)
	Extension Extension

	// BeginTime is the start time filter (optional, empty string to skip)
	// Format: "2006-01-02", only accept this format
	BeginTime string

	// EndTime is the end time filter (optional, empty string to skip)
	// Format: "2006-01-03", only accept this format
	EndTime string

	// MaxConcurrency is the maximum number of concurrent downloads
	// Only used in DownloadAllActivitiesWithConcurrency
	// Default: 5 (if set to 0)
	MaxConcurrency int

	// Callback is called for each downloaded activity
	// Return true to continue, false to stop downloading
	Callback DownloadCallback
}
