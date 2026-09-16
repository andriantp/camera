package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type ptzAbsoluteExXML struct {
	Elevation    float64 `xml:"elevation"`
	Azimuth      float64 `xml:"azimuth"`
	AbsoluteZoom float64 `xml:"absoluteZoom"`
	Focus        float64 `xml:"focus"`
}

func ParseAbsolute(resp *http.Response) (*model.Absolute, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, body)
	}

	var src ptzAbsoluteExXML

	if err := xml.Unmarshal(body, &src); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	return &model.Absolute{
		Position: model.PanTilt{
			X: src.Azimuth * 10,
			Y: src.Elevation * 10,
		},
		Zoom:  src.AbsoluteZoom * 10,
		Focus: src.Focus,
	}, nil
}
