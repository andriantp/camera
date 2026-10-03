package model

type ThermometryRulesTemperatureInfo struct {
	MaxTemperature      float64          `json:"max_temperature"`
	MinTemperature      float64          `json:"min_temperature"`
	AverageTemperature  float64          `json:"average_temperature"`
	MaxTemperaturePoint ThermometryPoint `json:"max_temperature_point"`
	MinTemperaturePoint ThermometryPoint `json:"min_temperature_point"`
	IsFreezeData        bool             `json:"is_freeze_data"`
}

type ThermometryPoint struct {
	PositionX float64 `json:"position_x"`
	PositionY float64 `json:"position_y"`
}
