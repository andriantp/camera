package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type envelopeProfileXML struct {
	XMLName xml.Name       `xml:"Envelope"`
	Header  struct{}       `xml:"Header"`
	Body    bodyProfileXML `xml:"Body"`
}

type bodyProfileXML struct {
	Response responseProfileXML `xml:"GetProfilesResponse"`
}

type responseProfileXML struct {
	Profiles []profileXML `xml:"Profiles"`
}

type profileXML struct {
	Token                     string                       `xml:"token,attr"`
	Name                      string                       `xml:"Name"`
	VideoEncoderConfiguration videoEncoderConfigurationXML `xml:"VideoEncoderConfiguration"`
}

type videoEncoderConfigurationXML struct {
	Encoding   string        `xml:"Encoding"`
	Resolution resolutionXML `xml:"Resolution"`
}

type resolutionXML struct {
	Width  int `xml:"Width"`
	Height int `xml:"Height"`
}

func ParseProfiles(resp *http.Response) ([]model.Profile, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}

	var env envelopeProfileXML

	if err := xml.Unmarshal(body, &env); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	if len(env.Body.Response.Profiles) == 0 {
		return nil, fmt.Errorf("no profile found")
	}

	profiles := make([]model.Profile, 0, len(env.Body.Response.Profiles))

	for _, p := range env.Body.Response.Profiles {

		profiles = append(profiles, model.Profile{
			Token: p.Token,
			Name:  p.Name,
			VideoEncoderConfiguration: model.VideoEncoderConfiguration{
				Encoding: p.VideoEncoderConfiguration.Encoding,
				Resolution: model.Resolution{
					Width:  p.VideoEncoderConfiguration.Resolution.Width,
					Height: p.VideoEncoderConfiguration.Resolution.Height,
				},
			},
		})
	}

	return profiles, nil
}
