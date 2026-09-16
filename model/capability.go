package model

type Capabilities struct {
	Onvif *OnvifCapabilities `json:"onvif,omitempty"`
	Isapi *IsapiCapabilities `json:"isapi,omitempty"`
}

type OnvifCapabilities struct {
	Analytics Endpoint `json:"analytics"`
	Device    Endpoint `json:"device"`
	Events    Endpoint `json:"events"`
	Imaging   Endpoint `json:"imaging"`
	PTZ       Endpoint `json:"ptz"`
	Media     Media    `json:"media"`
}

type IsapiCapabilities struct {
	System    *SystemCapabilities    `json:"system,omitempty"`
	Streaming *StreamingCapabilities `json:"streaming,omitempty"`
	PTZ       *PTZCapabilities       `json:"ptz,omitempty"`
	Event     *EventCapability       `json:"event,omitempty"`
	Image     *ImageCapabilities     `json:"image,omitempty"`
	Analytic  *AnalyticsCapability   `json:"analytic,omitempty"`
}

type Endpoint struct {
	XAddr string `json:"xaddr"`
}

type Media struct {
	XAddr     string    `json:"xaddr"`
	Streaming Streaming `json:"streaming"`
}

type Streaming struct {
	RTPMulticast bool `json:"rtp_multicast"`
	RTPTCP       bool `json:"rtp_tcp"`
	RTPRTSPTCP   bool `json:"rtp_rtsp_tcp"`
}
