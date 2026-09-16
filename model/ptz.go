package model

type PTZCapabilities struct {
	AbsolutePanTilt   *Range `json:"absolute_pan_tilt,omitempty"`
	RelativePanTilt   *Range `json:"relative_pan_tilt,omitempty"`
	ContinuousPanTilt *Range `json:"continuous_pan_tilt,omitempty"`
	MomentaryPanTilt  *Range `json:"momentary_pan_tilt,omitempty"`

	Zoom           CapabilitiesZoom   `json:"zoom"`
	Preset         CapabilitiesPreset `json:"preset"`
	ContinuousMove bool               `json:"continuous_move"`
	Position3D     bool               `json:"position_3d"`
}

type CapabilitiesPreset struct {
	Support   bool `json:"support"`
	MaxPreset int  `json:"max_preset"`
}

type CapabilitiesZoom struct {
	Absolute   AxisRange `json:"absolute,omitempty"`
	Continuous AxisRange `json:"continuous,omitempty"`
}
