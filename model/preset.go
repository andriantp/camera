package model

type Preset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

	Position Position `json:"position"`
}

type Position struct {
	Elevation int `json:"elevation"`
	Azimuth   int `json:"azimuth"`
	Zoom      int `json:"zoom"`
}
