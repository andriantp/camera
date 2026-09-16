package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type presetListXML struct {
	Presets []presetXML `xml:"PTZPreset"`
}

type presetXML struct {
	ID         int    `xml:"id"`
	Enabled    bool   `xml:"enabled"`
	PresetName string `xml:"presetName"`

	AbsoluteHigh absoluteHighXML `xml:"AbsoluteHigh"`
}

type absoluteHighXML struct {
	Elevation int `xml:"elevation"`
	Azimuth   int `xml:"azimuth"`
	Zoom      int `xml:"absoluteZoom"`
}

func ParsePresets(resp *http.Response) ([]model.Preset, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}

	var src presetListXML
	if err := xml.Unmarshal(body, &src); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	presets := make([]model.Preset, 0, len(src.Presets))

	for _, p := range src.Presets {
		presets = append(presets, model.Preset{
			ID:   p.ID,
			Name: p.PresetName,
			Position: model.Position{
				Elevation: p.AbsoluteHigh.Elevation,
				Azimuth:   p.AbsoluteHigh.Azimuth,
				Zoom:      p.AbsoluteHigh.Zoom,
			},
		})
	}

	return presets, nil
}
