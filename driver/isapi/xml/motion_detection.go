package xml

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"regexp"

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

func ParseMotionDetection(resp *http.Response) (*model.MotionDetection, error) {
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

func SetVMDSensitivity(body []byte, level int) ([]byte, error) {
	if level < 0 || level > 100 {
		return nil, fmt.Errorf(
			"VMD sensitivity must be between 0 and 100",
		)
	}

	re := regexp.MustCompile(
		`(<MotionDetectionLayout\b[^>]*>[\s\S]*?<sensitivityLevel>)[0-9]+(</sensitivityLevel>)`,
	)

	if !re.Match(body) {
		return nil, fmt.Errorf("sensitivityLevel not found")
	}

	replacement := []byte(fmt.Sprintf("${1}%d${2}", level))
	return re.ReplaceAll(body, replacement), nil
}

type responseStatus struct {
	StatusCode    int    `xml:"statusCode"`
	StatusString  string `xml:"statusString"`
	SubStatusCode string `xml:"subStatusCode"`
}

func ParseResponseStatus(resp *http.Response) error {
	var data responseStatus

	if err := xml.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("decode response status: %w", err)
	}

	if data.StatusCode != 1 || data.SubStatusCode != "ok" {
		return fmt.Errorf(
			"camera rejected request: status %d, %s (%s)",
			data.StatusCode,
			data.StatusString,
			data.SubStatusCode,
		)
	}

	return nil
}
