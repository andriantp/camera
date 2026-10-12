package model

type MotionDetection struct {
	Enabled          bool                  `json:"enabled"`
	EnableHighlight  bool                  `json:"enable_highlight"`
	SamplingInterval int                   `json:"sampling_interval"`
	StartTriggerTime int                   `json:"start_trigger_time"`
	EndTriggerTime   int                   `json:"end_trigger_time"`
	RegionType       string                `json:"region_type"`
	Grid             MotionDetectionGrid   `json:"grid"`
	Layout           MotionDetectionLayout `json:"layout"`
}

type MotionDetectionGrid struct {
	RowGranularity    int `json:"row_granularity"`
	ColumnGranularity int `json:"column_granularity"`
}

type MotionDetectionLayout struct {
	SensitivityLevel int    `json:"sensitivity_level"`
	GridMap          string `json:"grid_map"`
}
