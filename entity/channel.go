package entity

type (
	AvgMobility struct {
		Inside float64 `json:"inside"`
		Enter  float64 `json:"enter"`
		Exit   float64 `json:"exit"`
	}
)
