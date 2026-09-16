package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type ptzStatusXML struct {
	AbsoluteHigh statusabsoluteHighXML `xml:"AbsoluteHigh"`
}

type statusabsoluteHighXML struct {
	Elevation    float64 `xml:"elevation"`
	Azimuth      float64 `xml:"azimuth"`
	AbsoluteZoom float64 `xml:"absoluteZoom"`
}

func ParseStatus(resp *http.Response) (*model.Status, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}

	var src ptzStatusXML

	if err := xml.Unmarshal(body, &src); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	return &model.Status{
		Position: model.PanTilt{
			X: src.AbsoluteHigh.Azimuth,
			Y: src.AbsoluteHigh.Elevation,
		},
		Zoom: src.AbsoluteHigh.AbsoluteZoom,
	}, nil
}
