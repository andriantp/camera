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
