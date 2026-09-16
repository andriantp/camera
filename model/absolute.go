package model

type Absolute struct {
	Position PanTilt `json:"position"`
	Zoom     float64 `json:"zoom"`
	Focus    float64 `json:"focus"`
}
