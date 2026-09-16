package model

type Profile struct {
	Token string `json:"token"`
	Name  string `json:"name"`

	VideoEncoderConfiguration VideoEncoderConfiguration `json:"video_encoder_configuration"`
}

type VideoEncoderConfiguration struct {
	Encoding   string     `json:"encoding"`
	Resolution Resolution `json:"resolution"`
}
