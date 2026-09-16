package model

type PanTilt struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type AxisRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type Range struct {
	XRange AxisRange `json:"x_range"`
	YRange AxisRange `json:"y_range"`
}

type Resolution struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Speed struct {
	Pan  int `json:"pan"`
	Tilt int `json:"tilt"`
}


