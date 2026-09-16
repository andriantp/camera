package model

type SystemCapabilities struct {
	Network   NetworkCapabilities   `json:"network"`
	Video     VideoCapabilities     `json:"video"`
	Audio     AudioCapabilities     `json:"audio"`
	Snapshot  SnapshotCapabilities  `json:"snapshot"`
	VoiceTalk VoiceTalkCapabilities `json:"voice_talk"`
	Events    EventCapabilities     `json:"events"`
	Analytics AnalyticsCapabilities `json:"analytics"`
}

type NetworkCapabilities struct {
	NTP   bool `json:"ntp"`
	FTP   bool `json:"ftp"`
	UPnP  bool `json:"upnp"`
	DDNS  bool `json:"ddns"`
	HTTPS bool `json:"https"`
}

type VideoCapabilities struct {
	InputPorts  int  `json:"input_ports"`
	OutputPorts int  `json:"output_ports"`
	Heatmap     bool `json:"heatmap"`
	Counting    bool `json:"counting"`
	Picture     bool `json:"picture"`
	PrivacyMask bool `json:"privacy_mask"`
}

type AudioCapabilities struct {
	Input  int `json:"input"`
	Output int `json:"output"`
}

type SnapshotCapabilities struct {
	Supported bool `json:"supported"`
}

type VoiceTalkCapabilities struct {
	Channels int `json:"channels"`
}

type EventCapabilities struct {
	Subscribe bool `json:"subscribe"`
}

type AnalyticsCapabilities struct {
	ROI            bool `json:"roi"`
	AudioDetection bool `json:"audio_detection"`
	FaceDetection  bool `json:"face_detection"`
	LineDetection  bool `json:"line_detection"`
	FieldDetection bool `json:"field_detection"`
	RegionEntrance bool `json:"region_entrance"`
	RegionExiting  bool `json:"region_exiting"`
}
