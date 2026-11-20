package xml

import (
	"math"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestParseAbsoluteScalesToTenthsOfDegree(t *testing.T) {
	resp := response(http.StatusOK, `<PTZAbsoluteEx>
	<elevation>45.5</elevation>
	<azimuth>180</azimuth>
	<absoluteZoom>1.5</absoluteZoom>
	<focus>3</focus>
</PTZAbsoluteEx>`)

	got, err := ParseAbsolute(resp)
	if err != nil {
		t.Fatal(err)
	}

	if got.Position.X != 1800 || got.Position.Y != 455 || math.Abs(got.Zoom-15) > 1e-9 || got.Focus != 3 {
		t.Errorf("got %+v", got)
	}
}

func TestParsePTZCapabilities(t *testing.T) {
	resp := response(http.StatusOK, `<PTZChanelCap>
	<AbsolutePanTiltPositionSpace>
		<XRange><Min>0</Min><Max>3600</Max></XRange>
		<YRange><Min>-900</Min><Max>2700</Max></YRange>
	</AbsolutePanTiltPositionSpace>
	<AbsoluteZoomPositionSpace><ZRange><Min>10</Min><Max>320</Max></ZRange></AbsoluteZoomPositionSpace>
	<ContinuousPanTiltSpace>
		<XRange><Min>-100</Min><Max>100</Max></XRange>
		<YRange><Min>-100</Min><Max>100</Max></YRange>
	</ContinuousPanTiltSpace>
	<maxPresetNum>256</maxPresetNum>
	<PresetNameCap><presetNameSupport>true</presetNameSupport></PresetNameCap>
</PTZChanelCap>`)

	got, err := ParsePTZCapabilities(resp)
	if err != nil {
		t.Fatal(err)
	}

	if got.AbsolutePanTilt == nil || got.AbsolutePanTilt.XRange.Max != 3600 || got.AbsolutePanTilt.YRange.Min != -900 {
		t.Errorf("AbsolutePanTilt = %+v", got.AbsolutePanTilt)
	}
	if got.Zoom.Absolute.Min != 10 || got.Zoom.Absolute.Max != 320 {
		t.Errorf("Zoom = %+v", got.Zoom)
	}
	if !got.ContinuousMove || got.ContinuousPanTilt == nil {
		t.Error("continuous pan/tilt not detected")
	}
	if got.RelativePanTilt != nil || got.MomentaryPanTilt != nil {
		t.Error("missing spaces should be nil")
	}
	if !got.Preset.Support || got.Preset.MaxPreset != 256 {
		t.Errorf("Preset = %+v", got.Preset)
	}
}

func TestParseErrorsIncludeStatusCode(t *testing.T) {
	parsers := map[string]func(*http.Response) error{
		"ParseAbsolute":           func(r *http.Response) error { _, err := ParseAbsolute(r); return err },
		"ParsePTZCapabilities":    func(r *http.Response) error { _, err := ParsePTZCapabilities(r); return err },
		"ParseStatus":             func(r *http.Response) error { _, err := ParseStatus(r); return err },
		"ParseDeviceInfo":         func(r *http.Response) error { _, err := ParseDeviceInfo(r); return err },
		"ParsePresets":            func(r *http.Response) error { _, err := ParsePresets(r); return err },
		"ParseStreaming":          func(r *http.Response) error { _, err := ParseStreaming(r); return err },
		"ParseSystemCapabilities": func(r *http.Response) error { _, err := ParseSystemCapabilities(r); return err },
		"ParseStorage":            func(r *http.Response) error { _, err := ParseStorage(r); return err },
		"ParseContentSearch":      func(r *http.Response) error { _, err := ParseContentSearch(r); return err },
	}

	for name, parse := range parsers {
		t.Run(name, func(t *testing.T) {
			resp, body := trackedResponse(http.StatusUnauthorized, "Unauthorized")

			err := parse(resp)
			if err == nil || !strings.Contains(err.Error(), "401") {
				t.Errorf("err = %v, want it to contain status 401", err)
			}
			if !body.closed {
				t.Error("body was not closed")
			}
		})
	}
}

func TestParsersCloseBody(t *testing.T) {
	parsers := map[string]struct {
		body  string
		parse func(*http.Response) error
	}{
		"ParseMotionDetection": {
			`<MotionDetection><enabled>true</enabled></MotionDetection>`,
			func(r *http.Response) error { _, err := ParseMotionDetection(r); return err },
		},
		"ParseHTTPHosts": {
			`<HttpHostNotificationList></HttpHostNotificationList>`,
			func(r *http.Response) error { _, err := ParseHTTPHosts(r); return err },
		},
		"ParseResponseStatus": {
			`<ResponseStatus><statusCode>1</statusCode><subStatusCode>ok</subStatusCode></ResponseStatus>`,
			ParseResponseStatus,
		},
	}

	for name, tt := range parsers {
		t.Run(name, func(t *testing.T) {
			resp, body := trackedResponse(http.StatusOK, tt.body)
			if err := tt.parse(resp); err != nil {
				t.Fatal(err)
			}
			if !body.closed {
				t.Error("body was not closed")
			}
		})
	}
}

func TestParseStorage(t *testing.T) {
	resp := response(http.StatusOK, `<hddList>
	<hdd><id>1</id><hddName>hdd1</hddName><status>ok</status><capacity>1000</capacity><freeSpace>250</freeSpace></hdd>
</hddList>`)

	got, err := ParseStorage(resp)
	if err != nil {
		t.Fatal(err)
	}

	if got.Total != 1000 || got.Free != 250 || got.Used != 750 || got.UsedPct != 75 {
		t.Errorf("got %+v", got)
	}
}

func TestParseStorageZeroCapacity(t *testing.T) {
	resp := response(http.StatusOK, `<hddList><hdd><capacity>0</capacity><freeSpace>0</freeSpace></hdd></hddList>`)

	got, err := ParseStorage(resp)
	if err != nil {
		t.Fatal(err)
	}
	if got.UsedPct != 0 {
		t.Errorf("UsedPct = %v, want 0", got.UsedPct)
	}
}

func TestParseContentSearch(t *testing.T) {
	resp := response(http.StatusOK, `<CMSearchResult>
	<searchID>abc</searchID>
	<responseStatus>true</responseStatus>
	<responseStatusStrg>MORE</responseStatusStrg>
	<totalMatches>10</totalMatches>
	<numOfMatches>2</numOfMatches>
	<matchList>
		<searchMatchItem>
			<trackID>101</trackID>
			<timeSpan><startTime>2026-01-01T00:00:00Z</startTime><endTime>2026-01-01T01:01:01Z</endTime></timeSpan>
			<mediaSegmentDescriptor>
				<contentType>video</contentType>
				<codecType>H.264</codecType>
				<playbackURI>rtsp://cam/Streaming/tracks/101?starttime=20260101T000000Z&amp;size=2097152</playbackURI>
			</mediaSegmentDescriptor>
		</searchMatchItem>
		<searchMatchItem>
			<trackID>101</trackID>
			<timeSpan><startTime>2026-01-01T02:00:00Z</startTime><endTime>2026-01-01T02:00:30Z</endTime></timeSpan>
			<mediaSegmentDescriptor>
				<playbackURI>rtsp://cam/Streaming/tracks/101?size=1024</playbackURI>
			</mediaSegmentDescriptor>
		</searchMatchItem>
	</matchList>
</CMSearchResult>`)

	got, err := ParseContentSearch(resp)
	if err != nil {
		t.Fatal(err)
	}

	if got.SearchID != "abc" || !got.More || got.TotalMatches != 10 || len(got.Matches) != 2 {
		t.Fatalf("got %+v", got)
	}

	first := got.Matches[0]
	if first.Size.Bytes != 2097152 || first.Size.MB != 2 {
		t.Errorf("size = %+v", first.Size)
	}
	if first.Timeframe.Duration != "1h 1m 1s" {
		t.Errorf("duration = %q", first.Timeframe.Duration)
	}
	if first.Filename != "20260101T000000Z_20260101T010101Z.mp4" {
		t.Errorf("filename = %q", first.Filename)
	}

	// Each record must point at its own playback URI.
	if *got.Matches[0].PlaybackURI == *got.Matches[1].PlaybackURI {
		t.Error("records share the same PlaybackURI pointer value")
	}
}

func TestParseContentSearchMissingSize(t *testing.T) {
	resp := response(http.StatusOK, `<CMSearchResult><matchList><searchMatchItem>
		<timeSpan><startTime>2026-01-01T00:00:00Z</startTime><endTime>2026-01-01T00:00:01Z</endTime></timeSpan>
		<mediaSegmentDescriptor><playbackURI>rtsp://cam/track</playbackURI></mediaSegmentDescriptor>
	</searchMatchItem></matchList></CMSearchResult>`)

	if _, err := ParseContentSearch(resp); err == nil {
		t.Fatal("expected error when playback URI has no size")
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		in   time.Duration
		want string
	}{
		{0, "0s"},
		{30 * time.Second, "30s"},
		{90 * time.Second, "1m 30s"},
		{time.Hour, "1h"},
		{time.Hour + time.Minute + time.Second, "1h 1m 1s"},
		{1500 * time.Millisecond, "2s"},
	}

	for _, tt := range tests {
		if got := formatDuration(tt.in); got != tt.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParsePlaybackURI(t *testing.T) {
	got, err := parsePlaybackURI("rtsp://cam/track?starttime=x&size=12345")
	if err != nil || got != 12345 {
		t.Errorf("got %d, %v; want 12345", got, err)
	}

	for _, bad := range []string{"rtsp://cam/track", "rtsp://cam/track?size=abc"} {
		if _, err := parsePlaybackURI(bad); err == nil {
			t.Errorf("parsePlaybackURI(%q) should fail", bad)
		}
	}
}
