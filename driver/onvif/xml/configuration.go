package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type envelopeConfigurationXML struct {
	XMLName xml.Name             `xml:"Envelope"`
	Header  struct{}             `xml:"Header"`
	Body    bodyConfigurationXML `xml:"Body"`
}

type bodyConfigurationXML struct {
	Response responseConfigurationXML `xml:"GetConfigurationsResponse"`
}

type responseConfigurationXML struct {
	PTZConfiguration configurationXML `xml:"PTZConfiguration"`
}

type configurationXML struct {
	DefaultPTZSpeed defaultPTZSpeedXML `xml:"DefaultPTZSpeed"`
	PanTiltLimits   panTiltLimitsXML   `xml:"PanTiltLimits"`
}

type defaultPTZSpeedXML struct {
	PanTilt panTiltXML `xml:"PanTilt"`
}

type panTiltXML struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
}

type panTiltLimitsXML struct {
	Range rangeXML `xml:"Range"`
}

type rangeXML struct {
	XRange axisRangeXML `xml:"XRange"`
	YRange axisRangeXML `xml:"YRange"`
}

type axisRangeXML struct {
	Min float64 `xml:"Min"`
	Max float64 `xml:"Max"`
}

func ParseConfiguration(resp *http.Response) (*model.Configuration, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}

	var env envelopeConfigurationXML

	if err := xml.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	src := env.Body.Response.PTZConfiguration

	return &model.Configuration{
		DefaultPTZSpeed: model.DefaultPTZSpeed{
			PanTilt: model.PanTilt{
				X: src.DefaultPTZSpeed.PanTilt.X,
				Y: src.DefaultPTZSpeed.PanTilt.Y,
			},
		},
		PanTiltLimits: model.PanTiltLimits{
			Range: model.Range{
				XRange: model.AxisRange{
					Min: src.PanTiltLimits.Range.XRange.Min,
					Max: src.PanTiltLimits.Range.XRange.Max,
				},
				YRange: model.AxisRange{
					Min: src.PanTiltLimits.Range.YRange.Min,
					Max: src.PanTiltLimits.Range.YRange.Max,
				},
			},
		},
	}, nil
}
