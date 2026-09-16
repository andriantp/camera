package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type envelopeStatusXML struct {
	XMLName xml.Name      `xml:"Envelope"`
	Header  struct{}      `xml:"Header"`
	Body    bodyStatusXML `xml:"Body"`
}

type bodyStatusXML struct {
	Response responseStatusXML `xml:"GetStatusResponse"`
}

type responseStatusXML struct {
	PTZStatus ptzStatusXML `xml:"PTZStatus"`
}

type ptzStatusXML struct {
	Position   positionXML   `xml:"Position"`
	MoveStatus moveStatusXML `xml:"MoveStatus"`
	Error      int           `xml:"Error"`
}

type positionXML struct {
	PanTilt panTiltXML `xml:"PanTilt"`
}

type moveStatusXML struct {
	PanTilt string `xml:"PanTilt"`
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

	var env envelopeStatusXML

	if err := xml.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	src := env.Body.Response.PTZStatus

	return &model.Status{
		Position: model.PanTilt{
			X: src.Position.PanTilt.X,
			Y: src.Position.PanTilt.Y,
		},
		MoveStatus: src.MoveStatus.PanTilt,
		Error:      src.Error,
	}, nil
}
