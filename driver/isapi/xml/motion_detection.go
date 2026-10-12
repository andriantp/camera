package xml

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/andriantp/camera/model"
)

type motionDetection struct {
	Enabled          bool                  `xml:"enabled"`
	EnableHighlight  bool                  `xml:"enableHighlight"`
	SamplingInterval int                   `xml:"samplingInterval"`
	StartTriggerTime int                   `xml:"startTriggerTime"`
	EndTriggerTime   int                   `xml:"endTriggerTime"`
	RegionType       string                `xml:"regionType"`
	Grid             motionDetectionGrid   `xml:"Grid"`
	Layout           motionDetectionLayout `xml:"MotionDetectionLayout"`
}

type motionDetectionGrid struct {
	RowGranularity    int `xml:"rowGranularity"`
	ColumnGranularity int `xml:"columnGranularity"`
}

type motionDetectionLayout struct {
	SensitivityLevel int                    `xml:"sensitivityLevel"`
	Layout           motionDetectionGridMap `xml:"layout"`
}

type motionDetectionGridMap struct {
	GridMap string `xml:"gridMap"`
}

func ParseMotionDetection(
	resp *http.Response,
) (*model.MotionDetection, error) {
	var data motionDetection

	if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("decode motion detection: %w", err)
	}

	return &model.MotionDetection{
		Enabled:          data.Enabled,
		EnableHighlight:  data.EnableHighlight,
		SamplingInterval: data.SamplingInterval,
		StartTriggerTime: data.StartTriggerTime,
		EndTriggerTime:   data.EndTriggerTime,
		RegionType:       data.RegionType,
		Grid: model.MotionDetectionGrid{
			RowGranularity:    data.Grid.RowGranularity,
			ColumnGranularity: data.Grid.ColumnGranularity,
		},
		Layout: model.MotionDetectionLayout{
			SensitivityLevel: data.Layout.SensitivityLevel,
			GridMap:          data.Layout.Layout.GridMap,
		},
	}, nil
}
