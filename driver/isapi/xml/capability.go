package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type ptzCapabilitiesXML struct {
	AbsolutePanTilt   *panTiltSpaceXML `xml:"AbsolutePanTiltPositionSpace"`
	AbsoluteZoom      *zoomSpaceXML    `xml:"AbsoluteZoomPositionSpace"`
	ContinuousZoom    *zoomSpaceXML    `xml:"ContinuousZoomSpace"`
	RelativePanTilt   *panTiltSpaceXML `xml:"RelativePanTiltSpace"`
	ContinuousPanTilt *panTiltSpaceXML `xml:"ContinuousPanTiltSpace"`
	MomentaryPanTilt  *panTiltSpaceXML `xml:"MomentaryPanTiltSpace"`

	MaxPresetNum      int              `xml:"maxPresetNum"`
	Position3DSupport bool             `xml:"isSupportPosition3D"`
	PresetNameCap     presetNameCapXML `xml:"PresetNameCap"`
}

type presetNameCapXML struct {
	PresetNameSupport bool `xml:"presetNameSupport"`
}

type panTiltSpaceXML struct {
	XRange axisRangeXML `xml:"XRange"`
	YRange axisRangeXML `xml:"YRange"`
}

type zoomSpaceXML struct {
	ZRange axisRangeXML `xml:"ZRange"`
}

type axisRangeXML struct {
	Min float64 `xml:"Min"`
	Max float64 `xml:"Max"`
}

func ParsePTZCapabilities(resp *http.Response) (*model.PTZCapabilities, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}

	var src ptzCapabilitiesXML
	if err := xml.Unmarshal(body, &src); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	capabilities := &model.PTZCapabilities{
		Preset: model.CapabilitiesPreset{
			Support:   src.PresetNameCap.PresetNameSupport,
			MaxPreset: src.MaxPresetNum,
		},
		ContinuousMove: src.ContinuousPanTilt != nil,
		Position3D:     src.Position3DSupport,
	}

	if src.AbsoluteZoom != nil {
		capabilities.Zoom.Absolute = model.AxisRange{
			Min: src.AbsoluteZoom.ZRange.Min,
			Max: src.AbsoluteZoom.ZRange.Max,
		}
	}

	if src.ContinuousZoom != nil {
		capabilities.Zoom.Continuous = model.AxisRange{
			Min: src.ContinuousZoom.ZRange.Min,
			Max: src.ContinuousZoom.ZRange.Max,
		}
	}

	if src.AbsolutePanTilt != nil {
		capabilities.AbsolutePanTilt = toRange(src.AbsolutePanTilt)
	}

	if src.RelativePanTilt != nil {
		capabilities.RelativePanTilt = toRange(src.RelativePanTilt)
	}

	if src.ContinuousPanTilt != nil {
		capabilities.ContinuousPanTilt = toRange(src.ContinuousPanTilt)
	}

	if src.MomentaryPanTilt != nil {
		capabilities.MomentaryPanTilt = toRange(src.MomentaryPanTilt)
	}

	return capabilities, nil
}

func toRange(src *panTiltSpaceXML) *model.Range {
	return &model.Range{
		XRange: model.AxisRange{
			Min: src.XRange.Min,
			Max: src.XRange.Max,
		},
		YRange: model.AxisRange{
			Min: src.YRange.Min,
			Max: src.YRange.Max,
		},
	}
}
