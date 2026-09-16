package model

type Configuration struct {
	DefaultPTZSpeed DefaultPTZSpeed `json:"default_ptz_speed"`
	PanTiltLimits   PanTiltLimits   `json:"pan_tilt_limits"`
}

type DefaultPTZSpeed struct {
	PanTilt PanTilt `json:"pan_tilt"`
}

type PanTiltLimits struct {
	Range Range `json:"range"`
}
