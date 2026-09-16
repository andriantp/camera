package model

type Status struct {
	Position   PanTilt `json:"position"`
	Zoom       float64 `json:"zoom,omitempty"`
	MoveStatus string  `json:"move_status,omitempty"`
	Error      int     `json:"error,omitempty"`
}
