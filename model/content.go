package model

import "time"

type ContentCapabilities struct {
	NASNums            int      `json:"nas_nums"`
	PictureSearchTypes []string `json:"picture_search_types,omitempty"`
	RecordSearchTypes  []string `json:"record_search_types,omitempty"`
}

type Storage struct {
	Total   float64 `json:"total"`
	Used    float64 `json:"used"`
	Free    float64 `json:"free"`
	Unit    string  `json:"unit"`
	UsedPct float64 `json:"used_pct"`
}

type SearchResult struct {
	SearchID     string   `json:"search_id"`
	More         bool     `json:"more"`
	TotalMatches int      `json:"total_matches"`
	NumOfMatches int      `json:"num_of_matches"`
	Matches      []Record `json:"matches"`
}

type Timeframe struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration string    `json:"duration"`
}

type Size struct {
	Bytes int64   `json:"bytes"`
	MB    float64 `json:"mb"`
}

type Type struct {
	Content string `json:"content"`
	Codec   string `json:"codec"`
}

type Record struct {
	TrackID     int       `json:"track_id"`
	Filename    string    `json:"filename"`
	Timeframe   Timeframe `json:"timeframe"`
	Size        Size      `json:"size"`
	Type        Type      `json:"type"`
	PlaybackURI *string   `json:"playback_uri,omitempty"`
}
