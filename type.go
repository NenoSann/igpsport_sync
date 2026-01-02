package igpsportsync

type Extension string

const (
	FIT Extension = "0"
	GPX Extension = "1"
	TCX Extension = "2"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ActivityRow struct {
	RideID       int     `json:"rideId"`
	Title        string  `json:"title"`
	RideDistance float64 `json:"rideDistance"`
	StartTime    string  `json:"startTime"`
	ProductName  string  `json:"productName"`
}

type ActivityListData struct {
	Rows      []ActivityRow `json:"rows"`
	TotalPage int           `json:"totalPage"`
	PageNo    int           `json:"pageNo"`
	PageSize  int           `json:"pageSize"`
	TotalRows int           `json:"totalRows"`
}

type ActivityListResponse struct {
	Response
	Data ActivityListData `json:"data"`
}

type LoginResponse struct {
	Response
	Data LoginResult
}

type UserInfoResponse struct {
	Response
	Data UserInfoResult `json:"data"`
}

type ActivityDetailResponse struct {
	Response
	Data ActivityDetailData `json:"data"`
}

type LoginResult struct {
	Token_type    string `json:"token_type"`
	Access_token  string `json:"access_token"`
	Refresh_token string `json:"refresh_token"`
	Expires_in    int    `json:"expires_in"`
	Scope         string `json:"scope"`
	BoundPhone    bool   `json:"boundPhone"`
}

type UserInfoResult struct {
	StrMemberId     string  `json:"strMemberId"`
	MemberId        int     `json:"memberId"`
	NickName        string  `json:"nickName"`
	Avatar          string  `json:"avatar"`
	Sex             int     `json:"sex"`
	CityId          int     `json:"cityId"`
	CityName        string  `json:"cityName"`
	ProvinceName    string  `json:"provinceName"`
	Height          int     `json:"height"`
	Weight          float64 `json:"weight"`
	BirthDate       string  `json:"birthDate"`
	Mhr             int     `json:"mhr"`
	Lthr            int     `json:"lthr"`
	RideTime        int     `json:"rideTime"`
	RideDistance    int     `json:"rideDistance"`
	RideCalorie     int     `json:"rideCalorie"`
	RideNum         int     `json:"rideNum"`
	TimeZone        int     `json:"timeZone"`
	BikeWeight      float64 `json:"bikeWeight"`
	BikeWheelSize   float64 `json:"bikeWheelSize"`
	RegTime         string  `json:"regTime"`
	Integral        int     `json:"integral"`
	VO2max          int     `json:"vO2max"`
	Ftp             int     `json:"ftp"`
	Attention       int     `json:"attention"`
	Fans            int     `json:"fans"`
	UnitMetric      int     `json:"unitMetric"`
	UnitTemperature int     `json:"unitTemperature"`
	UnitWeight      int     `json:"unitWeight"`
	UnitHeight      int     `json:"unitHeight"`
	UnitLength      int     `json:"unitLength"`
	ViewFriends     int     `json:"viewFriends"`
	DeviceName      string  `json:"deviceName"`
	HasPassword     bool    `json:"hasPassword"`
	Phone           string  `json:"phone"`
	IsOfficial      bool    `json:"isOfficial"`
	MomentCount     int     `json:"momentCount"`
	IsIdentified    bool    `json:"isIdentified"`
	Type            int     `json:"type"`
	ShareUrl        string  `json:"shareUrl"`
}

type DeviceInfo struct {
	DeviceName      string `json:"deviceName"`
	DeviceImage     string `json:"deviceImage"`
	SoftwareVersion string `json:"softwareVersion"`
}

type ActivityDetailData struct {
	RideId               int        `json:"rideId"`
	MemberId             int        `json:"memberId"`
	Title                string     `json:"title"`
	Product              int        `json:"product"`
	SoftwareVersion      string     `json:"softwareVersion"`
	AvgSpeed             float64    `json:"avgSpeed"`
	AvgMovingSpeed       float64    `json:"avgMovingSpeed"`
	MaxSpeed             float64    `json:"maxSpeed"`
	StartTimeWithWeek    string     `json:"startTimeWithWeek"`
	StartTime            string     `json:"startTime"`
	MovingTime           int        `json:"movingTime"`
	EndTime              string     `json:"endTime"`
	TotalMovingTime      int        `json:"totalMovingTime"`
	TotalTime            int        `json:"totalTime"`
	RideDistance         int        `json:"rideDistance"`
	TotalAscent          int        `json:"totalAscent"`
	FitUrl               string     `json:"fitUrl"`
	Status               int        `json:"status"`
	ErrorType            int        `json:"errorType"`
	OpenStatus           int        `json:"openStatus"`
	DataSyncStravaStatus int        `json:"dataSyncStravaStatus"`
	DeviceInfo           DeviceInfo `json:"deviceInfo"`
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
