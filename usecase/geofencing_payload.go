package usecase

type (
	ResponseSSEDetectGeofence struct {
		Point       string `json:"point"`
		Detect      string `json:"detect"`
		ChannelName string `json:"channel_name"`
		DeviceId    string `json:"device_id"`
	}
)
