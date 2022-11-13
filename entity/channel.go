package entity

type (
	AvgMobility struct {
		Inside float64 `json:"inside"`
		Enter  float64 `json:"enter"`
		Exit   float64 `json:"exit"`
	}

	// data for send to websocket
	//get geofence by channel
	Counter struct {
		Enter       float64 `json:"enter"`
		Exit        float64 `json:"exit"`
		Inside      float64 `json:"inside"`
		CountDevice float64 `json:"count_device"`
	}
	DeviceType struct {
		Type     string `json:"type"`
		Count    int    `json:"count"`
		DateTime string `json:"datetime"`
	}
	PayloadMessageGetGeofenceByChann struct {
		Message     string       `json:"message"`
		ChannelName string       `json:"channel_name"`
		Counter     Counter      `json:"counter"`
		DeviceType  []DeviceType `json:"device_type"`
		Detect      string       `json:"detect"`
		Object      string       `json:"Object"`
	}
)
