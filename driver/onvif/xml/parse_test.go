package xml

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func response(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestParseProfiles(t *testing.T) {
	resp := response(http.StatusOK, `<Envelope><Body><GetProfilesResponse>
		<Profiles token="profile_1">
			<Name>mainStream</Name>
			<VideoEncoderConfiguration>
				<Encoding>H264</Encoding>
				<Resolution><Width>1920</Width><Height>1080</Height></Resolution>
			</VideoEncoderConfiguration>
		</Profiles>
		<Profiles token="profile_2"><Name>subStream</Name></Profiles>
	</GetProfilesResponse></Body></Envelope>`)

	got, err := ParseProfiles(resp)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("profiles = %d, want 2", len(got))
	}
	p := got[0]
	if p.Token != "profile_1" || p.Name != "mainStream" || p.VideoEncoderConfiguration.Encoding != "H264" ||
		p.VideoEncoderConfiguration.Resolution.Width != 1920 {
		t.Errorf("profile = %+v", p)
	}
}

func TestParseProfilesEmpty(t *testing.T) {
	resp := response(http.StatusOK, `<Envelope><Body><GetProfilesResponse></GetProfilesResponse></Body></Envelope>`)

	if _, err := ParseProfiles(resp); err == nil {
		t.Fatal("expected error for no profiles")
	}
}

func TestParseConfiguration(t *testing.T) {
	resp := response(http.StatusOK, `<Envelope><Body><GetConfigurationsResponse><PTZConfiguration>
		<DefaultPTZSpeed><PanTilt x="0.5" y="0.4"/></DefaultPTZSpeed>
		<PanTiltLimits><Range>
			<XRange><Min>-1</Min><Max>1</Max></XRange>
			<YRange><Min>-0.5</Min><Max>0.5</Max></YRange>
		</Range></PanTiltLimits>
	</PTZConfiguration></GetConfigurationsResponse></Body></Envelope>`)

	got, err := ParseConfiguration(resp)
	if err != nil {
		t.Fatal(err)
	}

	if got.DefaultPTZSpeed.PanTilt.X != 0.5 || got.PanTiltLimits.Range.XRange.Min != -1 || got.PanTiltLimits.Range.YRange.Max != 0.5 {
		t.Errorf("got %+v", got)
	}
}

func TestParseCapabilities(t *testing.T) {
	resp := response(http.StatusOK, `<Envelope><Body><GetCapabilitiesResponse><Capabilities>
		<Analytics><XAddr>http://cam/onvif/analytics</XAddr></Analytics>
		<Device><XAddr>http://cam/onvif/device</XAddr></Device>
		<Events><XAddr>http://cam/onvif/events</XAddr></Events>
		<Imaging><XAddr>http://cam/onvif/imaging</XAddr></Imaging>
		<PTZ><XAddr>http://cam/onvif/ptz</XAddr></PTZ>
		<Media>
			<XAddr>http://cam/onvif/media</XAddr>
			<StreamingCapabilities><RTPMulticast>false</RTPMulticast><RTP_TCP>true</RTP_TCP><RTP_RTSP_TCP>true</RTP_RTSP_TCP></StreamingCapabilities>
		</Media>
	</Capabilities></GetCapabilitiesResponse></Body></Envelope>`)

	got, err := ParseCapabilities(resp)
	if err != nil {
		t.Fatal(err)
	}

	o := got.Onvif
	if o.PTZ.XAddr != "http://cam/onvif/ptz" || o.Device.XAddr != "http://cam/onvif/device" ||
		o.Events.XAddr != "http://cam/onvif/events" || o.Imaging.XAddr != "http://cam/onvif/imaging" ||
		o.Analytics.XAddr != "http://cam/onvif/analytics" || o.Media.XAddr != "http://cam/onvif/media" {
		t.Errorf("endpoints = %+v", o)
	}
	if o.Media.Streaming.RTPMulticast || !o.Media.Streaming.RTPTCP || !o.Media.Streaming.RTPRTSPTCP {
		t.Errorf("streaming = %+v", o.Media.Streaming)
	}
}

func TestParseStatus(t *testing.T) {
	resp := response(http.StatusOK, `<Envelope><Body><GetStatusResponse><PTZStatus>
		<Position><PanTilt x="0.5" y="-0.25"/></Position>
		<MoveStatus><PanTilt>MOVING</PanTilt></MoveStatus>
		<Error>0</Error>
	</PTZStatus></GetStatusResponse></Body></Envelope>`)

	got, err := ParseStatus(resp)
	if err != nil {
		t.Fatal(err)
	}

	if got.Position.X != 0.5 || got.Position.Y != -0.25 || got.MoveStatus != "MOVING" || got.Error != 0 {
		t.Errorf("got %+v", got)
	}
}

func TestParsersRejectInvalidXML(t *testing.T) {
	parsers := map[string]func(*http.Response) error{
		"ParseCapabilities":  func(r *http.Response) error { _, err := ParseCapabilities(r); return err },
		"ParseConfiguration": func(r *http.Response) error { _, err := ParseConfiguration(r); return err },
		"ParseProfiles":      func(r *http.Response) error { _, err := ParseProfiles(r); return err },
		"ParseStatus":        func(r *http.Response) error { _, err := ParseStatus(r); return err },
	}

	for name, parse := range parsers {
		t.Run(name, func(t *testing.T) {
			if err := parse(response(http.StatusOK, "<broken")); err == nil {
				t.Error("expected error for invalid XML")
			}
		})
	}
}

func TestParseErrorsIncludeStatusCode(t *testing.T) {
	parsers := map[string]func(*http.Response) error{
		"ParseCapabilities":  func(r *http.Response) error { _, err := ParseCapabilities(r); return err },
		"ParseConfiguration": func(r *http.Response) error { _, err := ParseConfiguration(r); return err },
		"ParseProfiles":      func(r *http.Response) error { _, err := ParseProfiles(r); return err },
		"ParseStatus":        func(r *http.Response) error { _, err := ParseStatus(r); return err },
	}

	for name, parse := range parsers {
		t.Run(name, func(t *testing.T) {
			err := parse(response(http.StatusBadRequest, "SOAP fault"))
			if err == nil || !strings.Contains(err.Error(), "status 400") {
				t.Errorf("err = %v, want it to contain status 400", err)
			}
		})
	}
}
