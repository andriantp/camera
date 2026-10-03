package json

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andriantp/camera/model"
)

type thermometryRulesTemperatureInfoList struct {
	ThermometryRulesTemperatureInfo []thermometryRulesTemperatureInfo `json:"ThermometryRulesTemperatureInfo"`
}

type thermometryRulesTemperatureInfoResponse struct {
	ThermometryRulesTemperatureInfoList thermometryRulesTemperatureInfoList `json:"ThermometryRulesTemperatureInfoList"`
}

type thermometryRulesTemperatureInfo struct {
	MaxTemperature      float64          `json:"maxTemperature"`
	MinTemperature      float64          `json:"minTemperature"`
	AverageTemperature  float64          `json:"averageTemperature"`
	MaxTemperaturePoint thermometryPoint `json:"MaxTemperaturePoint"`
	MinTemperaturePoint thermometryPoint `json:"MinTemperaturePoint"`
	IsFreezeData        bool             `json:"isFreezedata"`
}

type thermometryPoint struct {
	PositionX float64 `json:"positionX"`
	PositionY float64 `json:"positionY"`
}

func ParseThermometryRulesTemperatureInfo(resp *http.Response) ([]model.ThermometryRulesTemperatureInfo, error) {
	var data thermometryRulesTemperatureInfoResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf(
			"decode thermometry rules temperature info: %w",
			err,
		)
	}

	result := make([]model.ThermometryRulesTemperatureInfo, 0)

	for _, item := range data.ThermometryRulesTemperatureInfoList.ThermometryRulesTemperatureInfo {
		result = append(result, model.ThermometryRulesTemperatureInfo{
			MaxTemperature:     item.MaxTemperature,
			MinTemperature:     item.MinTemperature,
			AverageTemperature: item.AverageTemperature,
			MaxTemperaturePoint: model.ThermometryPoint{
				PositionX: item.MaxTemperaturePoint.PositionX,
				PositionY: item.MaxTemperaturePoint.PositionY,
			},
			MinTemperaturePoint: model.ThermometryPoint{
				PositionX: item.MinTemperaturePoint.PositionX,
				PositionY: item.MinTemperaturePoint.PositionY,
			},
			IsFreezeData: item.IsFreezeData,
		})
	}

	return result, nil
}
