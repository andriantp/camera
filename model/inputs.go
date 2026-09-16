package model

type Input struct {
	ID          int    `json:"id"`
	InputPort   int    `json:"input_port"`
	Name        string `json:"name"`
	VideoFormat string `json:"video_format"`
}
