package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type envelopeCapabilities struct {
	XMLName xml.Name         `xml:"Envelope"`
	Header  struct{}         `xml:"Header"`
	Body    bodyCapabilities `xml:"Body"`
}

type bodyCapabilities struct {
	Response responseCapabilities `xml:"GetCapabilitiesResponse"`
}

type responseCapabilities struct {
	Capabilities capabilitiesXML `xml:"Capabilities"`
}

type capabilitiesXML struct {
	Analytics endpointXML `xml:"Analytics"`
	Device    endpointXML `xml:"Device"`
	Events    endpointXML `xml:"Events"`
	Imaging   endpointXML `xml:"Imaging"`
	PTZ       endpointXML `xml:"PTZ"`
	Media     mediaXML    `xml:"Media"`
}

type endpointXML struct {
	XAddr string `xml:"XAddr"`
}

type mediaXML struct {
	XAddr string `xml:"XAddr"`

	StreamingCapabilities streamingXML `xml:"StreamingCapabilities"`
}

type streamingXML struct {
	RTPMulticast bool `xml:"RTPMulticast"`
	RTPTCP       bool `xml:"RTP_TCP"`
	RTPRTSPTCP   bool `xml:"RTP_RTSP_TCP"`
}

func ParseCapabilities(resp *http.Response) (*model.Capabilities, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}

	var env envelopeCapabilities

	if err := xml.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	src := env.Body.Response.Capabilities

	return &model.Capabilities{
		Onvif: &model.OnvifCapabilities{
			Analytics: model.Endpoint{
				XAddr: src.Analytics.XAddr,
			},
			Device: model.Endpoint{
				XAddr: src.Device.XAddr,
			},
			Events: model.Endpoint{
				XAddr: src.Events.XAddr,
			},
			Imaging: model.Endpoint{
				XAddr: src.Imaging.XAddr,
			},
			PTZ: model.Endpoint{
				XAddr: src.PTZ.XAddr,
			},
			Media: model.Media{
				XAddr: src.Media.XAddr,
				Streaming: model.Streaming{
					RTPMulticast: src.Media.StreamingCapabilities.RTPMulticast,
					RTPTCP:       src.Media.StreamingCapabilities.RTPTCP,
					RTPRTSPTCP:   src.Media.StreamingCapabilities.RTPRTSPTCP,
				},
			},
		},
	}, nil
}
