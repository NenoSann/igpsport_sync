package igpsportsync

type Extension string

const (
	FIT Extension = "0"
	GPX Extension = "1"
	TCX Extension = "2"
)

type ActivityRow struct {
	RideID    int    `json:"rideId"`
	Title     string `json:"title"`
	StartTime string `json:"startTime"`
}

type ActivityListData struct {
	Rows      []ActivityRow `json:"rows"`
	TotalPage int           `json:"totalPage"`
}

type ActivityListResponse struct {
	Data ActivityListData `json:"data"`
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
