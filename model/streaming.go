package model

type StreamingChannel struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Video Video  `json:"video"`
	Audio Audio  `json:"audio"`
}

type Video struct {
	Codec      string     `json:"codec"`
	Resolution Resolution `json:"resolution"`
}

type Audio struct {
	Enabled bool `json:"enabled"`
}

type StreamingCapabilities struct {
	Transport TransportCapabilities      `json:"transport"`
	Video     StreamingVideoCapabilities `json:"video"`
	Audio     StreamingAudioCapabilities `json:"audio"`
	Features  StreamingFeatures          `json:"features"`
}

type TransportCapabilities struct {
	ControlProtocols []string     `json:"control_protocols"`
	Multicast        bool         `json:"multicast"`
	Unicast          bool         `json:"unicast"`
	RTPTransport     RTPTransport `json:"rtp_transport"`
	Security         bool         `json:"security"`
	CertificateTypes []string     `json:"certificate_types"`
}

type RTPTransport struct {
	Options []string `json:"options"`
	Type    string   `json:"type"`
}

type StreamingVideoCapabilities struct {
	Enabled       bool     `json:"enabled"`
	InputChannels []int    `json:"input_channels"`
	Resolutions   []string `json:"resolutions"`
	SVC           bool     `json:"svc"`
}

type StreamingAudioCapabilities struct {
	Enabled       bool  `json:"enabled"`
	InputChannels []int `json:"input_channels"`
}

type StreamingFeatures struct {
	DynamicCapability             bool `json:"dynamic_capability"`
	SmartCodeWithoutRestart       bool `json:"smart_code_without_restart"`
	RTCP                          bool `json:"rtcp"`
	SpecifiedImageQualitySnapshot bool `json:"specified_image_quality_snapshot"`
}
