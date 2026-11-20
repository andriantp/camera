package isapi

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/andriantp/camera/model"
)

type route struct {
	status      int
	contentType string
	body        string
}

type recordedRequest struct {
	method, path, query, body string
}

// routedCamera answers "METHOD /path" with a canned response and records
// every request it receives. Unknown routes answer 404.
type routedCamera struct {
	mu       sync.Mutex
	routes   map[string]route
	requests []recordedRequest
}

func newRoutedCamera(t *testing.T, routes map[string]route) (*Client, Camera, *routedCamera) {
	t.Helper()

	rc := &routedCamera{routes: routes}
	client, cam := newTestCamera(t, rc.serve)

	return client, cam, rc
}

func (rc *routedCamera) serve(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	rc.mu.Lock()
	rc.requests = append(rc.requests, recordedRequest{r.Method, r.URL.Path, r.URL.RawQuery, string(body)})
	rt, ok := rc.routes[r.Method+" "+r.URL.Path]
	rc.mu.Unlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if rt.contentType != "" {
		w.Header().Set("Content-Type", rt.contentType)
	}
	if rt.status != 0 {
		w.WriteHeader(rt.status)
	}
	_, _ = w.Write([]byte(rt.body))
}

// last returns the last recorded request with the given method.
func (rc *routedCamera) last(method string) recordedRequest {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	for i := len(rc.requests) - 1; i >= 0; i-- {
		if rc.requests[i].method == method {
			return rc.requests[i]
		}
	}
	return recordedRequest{}
}

const responseStatusOK = `<ResponseStatus><statusCode>1</statusCode><statusString>OK</statusString><subStatusCode>ok</subStatusCode></ResponseStatus>`

const responseStatusRejected = `<ResponseStatus><statusCode>4</statusCode><statusString>Invalid Operation</statusString><subStatusCode>notSupport</subStatusCode></ResponseStatus>`

func TestSimpleGetters(t *testing.T) {
	routes := map[string]route{
		"GET /ISAPI/PTZCtrl/channels/1/absoluteEx":                  {body: `<PTZAbsoluteEx><azimuth>90</azimuth><elevation>10</elevation><absoluteZoom>2</absoluteZoom></PTZAbsoluteEx>`},
		"GET /ISAPI/ContentMgmt/capabilities":                       {body: `<RacmCap><nasNums>8</nasNums></RacmCap>`},
		"GET /ISAPI/ContentMgmt/Storage/hdd":                        {body: `<hddList><hdd><capacity>100</capacity><freeSpace>40</freeSpace></hdd></hddList>`},
		"GET /ISAPI/System/deviceInfo":                              {body: `<DeviceInfo><model>DS-1</model></DeviceInfo>`},
		"GET /ISAPI/System/Video/inputs/channels":                   {body: `<VideoInputChannelList><VideoInputChannel><id>1</id></VideoInputChannel></VideoInputChannelList>`},
		"GET /ISAPI/PTZCtrl/channels/1/presets":                     {body: `<PTZPresetList><PTZPreset><id>3</id><presetName>Gate</presetName></PTZPreset></PTZPresetList>`},
		"GET /ISAPI/Streaming/channels":                             {body: `<StreamingChannelList><StreamingChannel><id>101</id></StreamingChannel></StreamingChannelList>`},
		"GET /ISAPI/Streaming/channels/1/capabilities":              {body: `<StreamingChannel><Video><enabled>true</enabled></Video></StreamingChannel>`},
		"GET /ISAPI/System/capabilities":                            {body: `<DeviceCap><isSupportSnapshot>true</isSupportSnapshot></DeviceCap>`},
		"GET /ISAPI/Event/capabilities":                             {body: `<EventCap><isSupportHDFull>true</isSupportHDFull></EventCap>`},
		"GET /ISAPI/Image/channels/1/capabilities":                  {body: `<ImageChannel><WDR><mode opt="open,close">open</mode></WDR></ImageChannel>`},
		"GET /ISAPI/Smart/Capabilities":                             {body: `<SmartCap><isSupportROI>true</isSupportROI></SmartCap>`},
		"GET /ISAPI/Event/notification/httpHosts":                   {body: `<HttpHostNotificationList><HttpHostNotification><id>1</id><ipAddress>10.0.0.5</ipAddress></HttpHostNotification></HttpHostNotificationList>`},
		"GET /ISAPI/System/Video/inputs/channels/1/motionDetection": {body: `<MotionDetection><enabled>true</enabled><MotionDetectionLayout><sensitivityLevel>40</sensitivityLevel></MotionDetectionLayout></MotionDetection>`},
	}
	client, cam, _ := newRoutedCamera(t, routes)
	ctx := context.Background()

	abs, err := client.GetAbsolute(ctx, cam)
	if err != nil || abs.Position.X != 900 || abs.Zoom != 20 {
		t.Errorf("GetAbsolute = %+v, %v", abs, err)
	}

	contentCaps, err := client.GetContentCapabilities(ctx, cam)
	if err != nil || contentCaps.NASNums != 8 {
		t.Errorf("GetContentCapabilities = %+v, %v", contentCaps, err)
	}

	storage, err := client.GetStorage(ctx, cam)
	if err != nil || storage.Used != 60 {
		t.Errorf("GetStorage = %+v, %v", storage, err)
	}

	info, err := client.GetDeviceInfo(ctx, cam)
	if err != nil || info.Model != "DS-1" {
		t.Errorf("GetDeviceInfo = %+v, %v", info, err)
	}

	inputs, err := client.GetInputs(ctx, cam)
	if err != nil || len(inputs) != 1 || inputs[0].ID != 1 {
		t.Errorf("GetInputs = %+v, %v", inputs, err)
	}

	presets, err := client.GetPresets(ctx, cam)
	if err != nil || len(presets) != 1 || presets[0].Name != "Gate" {
		t.Errorf("GetPresets = %+v, %v", presets, err)
	}

	streams, err := client.GetStreaming(ctx, cam)
	if err != nil || len(streams) != 1 || streams[0].ID != 101 {
		t.Errorf("GetStreaming = %+v, %v", streams, err)
	}

	streamCaps, err := client.GetStreamingCapabilities(ctx, cam)
	if err != nil || !streamCaps.Video.Enabled {
		t.Errorf("GetStreamingCapabilities = %+v, %v", streamCaps, err)
	}

	sysCaps, err := client.GetSystemCapabilities(ctx, cam)
	if err != nil || !sysCaps.Snapshot.Supported {
		t.Errorf("GetSystemCapabilities = %+v, %v", sysCaps, err)
	}

	eventCaps, err := client.GetEventCapabilities(ctx, cam)
	if err != nil || !eventCaps.Storage.HDFull {
		t.Errorf("GetEventCapabilities = %+v, %v", eventCaps, err)
	}

	imageCaps, err := client.GetImageCapabilities(ctx, cam)
	if err != nil || len(imageCaps.WDR.Modes) != 2 {
		t.Errorf("GetImageCapabilities = %+v, %v", imageCaps, err)
	}

	analytics, err := client.GetAnalyticsCapabilities(ctx, cam)
	if err != nil || !analytics.ROI {
		t.Errorf("GetAnalyticsCapabilities = %+v, %v", analytics, err)
	}

	hosts, err := client.GetHTTPHosts(ctx, cam)
	if err != nil || len(hosts) != 1 || hosts[0].IPAddress != "10.0.0.5" {
		t.Errorf("GetHTTPHosts = %+v, %v", hosts, err)
	}

	motion, err := client.GetMotionDetection(ctx, cam)
	if err != nil || !motion.Enabled || motion.Layout.SensitivityLevel != 40 {
		t.Errorf("GetMotionDetection = %+v, %v", motion, err)
	}
}

func TestGettersReturnStatusErrors(t *testing.T) {
	client, cam, _ := newRoutedCamera(t, map[string]route{}) // every route 404s
	ctx := context.Background()

	calls := map[string]func() error{
		"GetDeviceInfo":      func() error { _, err := client.GetDeviceInfo(ctx, cam); return err },
		"GetHTTPHosts":       func() error { _, err := client.GetHTTPHosts(ctx, cam); return err },
		"GetMotionDetection": func() error { _, err := client.GetMotionDetection(ctx, cam); return err },
		"GetThermometryRulesTemperatureInfo": func() error {
			_, err := client.GetThermometryRulesTemperatureInfo(ctx, cam, 1)
			return err
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			err := call()
			if err == nil || !strings.Contains(err.Error(), "404") {
				t.Errorf("err = %v, want status 404", err)
			}
		})
	}
}

func TestGettersReturnTransportErrors(t *testing.T) {
	client, cam := newTestCamera(t, func(w http.ResponseWriter, r *http.Request) {})
	cam.Host = "127.0.0.1:1" // nothing listens here

	ctx := context.Background()
	calls := map[string]func() error{
		"GetDeviceInfo": func() error { _, err := client.GetDeviceInfo(ctx, cam); return err },
		"GetSnapshot":   func() error { _, _, err := client.GetSnapshot(ctx, cam, SnapshotOptions{}); return err },
		"GotoHome":      func() error { return client.GotoHome(ctx, cam) },
		"Reboot":        func() error { return client.Reboot(ctx, cam) },
		"GotoPreset":    func() error { return client.GotoPreset(ctx, cam, 1) },
		"SetThermometry": func() error {
			return client.SetThermometry(ctx, cam, true)
		},
		"SetMotionSensitivity": func() error { return client.SetMotionSensitivity(ctx, cam, 50) },
		"ContentDownload": func() error {
			_, _, err := client.ContentDownload(ctx, cam, "rtsp://x")
			return err
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			if err := call(); err == nil {
				t.Error("expected connection error")
			}
		})
	}
}

func TestPTZCommands(t *testing.T) {
	client, cam, rc := newRoutedCamera(t, map[string]route{
		"PUT /ISAPI/PTZCtrl/channels/1/homeposition/goto": {},
		"PUT /ISAPI/PTZCtrl/channels/1/presets/7/goto":    {},
	})
	ctx := context.Background()

	if err := client.GotoHome(ctx, cam); err != nil {
		t.Errorf("GotoHome: %v", err)
	}

	if err := client.GotoPreset(ctx, cam, 7); err != nil {
		t.Errorf("GotoPreset: %v", err)
	}
	if req := rc.last(http.MethodPut); !strings.Contains(req.body, "<id>7</id>") {
		t.Errorf("GotoPreset body = %q", req.body)
	}

	// Unknown preset → camera answers 404.
	if err := client.GotoPreset(ctx, cam, 99); err == nil || !strings.Contains(err.Error(), "404") {
		t.Errorf("GotoPreset(99) err = %v, want 404", err)
	}
}

func TestPTZCommandsRejected(t *testing.T) {
	client, cam, _ := newRoutedCamera(t, map[string]route{})
	ctx := context.Background()

	if err := client.GotoHome(ctx, cam); err == nil {
		t.Error("GotoHome: expected error")
	}
}

func TestReboot(t *testing.T) {
	client, cam, rc := newRoutedCamera(t, map[string]route{
		"PUT /ISAPI/System/reboot": {body: responseStatusOK},
	})

	if err := client.Reboot(context.Background(), cam); err != nil {
		t.Fatal(err)
	}
	if req := rc.last(http.MethodPut); req.path != "/ISAPI/System/reboot" {
		t.Errorf("request = %+v", req)
	}

	client, cam, _ = newRoutedCamera(t, map[string]route{
		"PUT /ISAPI/System/reboot": {status: http.StatusForbidden, body: "forbidden"},
	})
	if err := client.Reboot(context.Background(), cam); err == nil || !strings.Contains(err.Error(), "403") {
		t.Errorf("err = %v, want status 403", err)
	}
}

func TestContentSearchRequest(t *testing.T) {
	client, cam, rc := newRoutedCamera(t, map[string]route{
		"POST /ISAPI/ContentMgmt/search": {body: `<CMSearchResult><searchID>s1</searchID><responseStatusStrg>OK</responseStatusStrg></CMSearchResult>`},
	})

	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := client.ContentSearch(context.Background(), cam, SearchRequest{
		SearchID:             1,
		TrackID:              101,
		StartTime:            start,
		EndTime:              start.Add(time.Hour),
		MaxResults:           40,
		SearchResultPosition: 80,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.SearchID != "s1" || result.More {
		t.Errorf("result = %+v", result)
	}

	body := rc.last(http.MethodPost).body
	for _, want := range []string{
		"<trackID>101</trackID>",
		"<startTime>2026-01-01T00:00:00Z</startTime>",
		"<endTime>2026-01-01T01:00:00Z</endTime>",
		"<maxResults>40</maxResults>",
		"<searchResultPosition>80</searchResultPosition>",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("request body missing %s:\n%s", want, body)
		}
	}
}

func TestContentDownloadRejected(t *testing.T) {
	client, cam, _ := newRoutedCamera(t, map[string]route{
		"POST /ISAPI/ContentMgmt/download": {status: http.StatusBadRequest, body: "bad playback uri"},
	})

	_, _, err := client.ContentDownload(context.Background(), cam, "rtsp://x")
	if err == nil || !strings.Contains(err.Error(), "bad playback uri") {
		t.Errorf("err = %v", err)
	}
}

func TestSetHTTPHost(t *testing.T) {
	const path = "PUT /ISAPI/Event/notification/httpHosts/1"

	host := model.HTTPHostNotification{
		ID:                       1,
		URL:                      "/events",
		ProtocolType:             "HTTP",
		ParameterFormatType:      "XML",
		AddressingFormatType:     "ipaddress",
		IPAddress:                "10.0.0.5",
		PortNo:                   8080,
		HTTPAuthenticationMethod: "none",
	}

	t.Run("enable", func(t *testing.T) {
		client, cam, rc := newRoutedCamera(t, map[string]route{path: {body: responseStatusOK}})

		if err := client.SetHTTPHost(context.Background(), cam, host, true); err != nil {
			t.Fatal(err)
		}
		body := rc.last(http.MethodPut).body
		if !strings.Contains(body, "<ipAddress>10.0.0.5</ipAddress>") || !strings.Contains(body, "<portNo>8080</portNo>") {
			t.Errorf("body = %s", body)
		}
	})

	t.Run("disable resets destination", func(t *testing.T) {
		client, cam, rc := newRoutedCamera(t, map[string]route{path: {body: responseStatusOK}})

		if err := client.SetHTTPHost(context.Background(), cam, host, false); err != nil {
			t.Fatal(err)
		}
		body := rc.last(http.MethodPut).body
		if !strings.Contains(body, "<ipAddress>0.0.0.0</ipAddress>") || !strings.Contains(body, "<portNo>80</portNo>") {
			t.Errorf("body = %s", body)
		}
	})

	t.Run("camera rejects", func(t *testing.T) {
		client, cam, _ := newRoutedCamera(t, map[string]route{path: {body: responseStatusRejected}})

		err := client.SetHTTPHost(context.Background(), cam, host, true)
		if err == nil || !strings.Contains(err.Error(), "Invalid Operation") {
			t.Errorf("err = %v", err)
		}
	})

	t.Run("http error", func(t *testing.T) {
		client, cam, _ := newRoutedCamera(t, map[string]route{path: {status: http.StatusBadRequest, body: "bad"}})

		if err := client.SetHTTPHost(context.Background(), cam, host, true); err == nil || !strings.Contains(err.Error(), "400") {
			t.Errorf("err = %v", err)
		}
	})

	t.Run("validation", func(t *testing.T) {
		client, cam, rc := newRoutedCamera(t, map[string]route{})

		invalidID := host
		invalidID.ID = 0
		if err := client.SetHTTPHost(context.Background(), cam, invalidID, true); err == nil {
			t.Error("expected error for ID 0")
		}

		missingIP := host
		missingIP.IPAddress = ""
		if err := client.SetHTTPHost(context.Background(), cam, missingIP, true); err == nil {
			t.Error("expected error for missing IP when enabling")
		}

		if len(rc.requests) != 0 {
			t.Errorf("validation errors should not reach the camera, got %d requests", len(rc.requests))
		}
	})
}

func TestSetMotionSensitivity(t *testing.T) {
	const path = "/ISAPI/System/Video/inputs/channels/1/motionDetection"
	const current = `<MotionDetection><enabled>true</enabled><MotionDetectionLayout><sensitivityLevel>20</sensitivityLevel></MotionDetectionLayout></MotionDetection>`

	t.Run("updates level and keeps other fields", func(t *testing.T) {
		client, cam, rc := newRoutedCamera(t, map[string]route{
			"GET " + path: {body: current},
			"PUT " + path: {body: responseStatusOK},
		})

		if err := client.SetMotionSensitivity(context.Background(), cam, 80); err != nil {
			t.Fatal(err)
		}
		body := rc.last(http.MethodPut).body
		if body != strings.Replace(current, ">20<", ">80<", 1) {
			t.Errorf("PUT body = %s", body)
		}
	})

	t.Run("invalid level", func(t *testing.T) {
		client, cam, rc := newRoutedCamera(t, map[string]route{})

		for _, level := range []int{-1, 101} {
			if err := client.SetMotionSensitivity(context.Background(), cam, level); err == nil {
				t.Errorf("level %d: expected error", level)
			}
		}
		if len(rc.requests) != 0 {
			t.Error("invalid level should not reach the camera")
		}
	})

	t.Run("get fails", func(t *testing.T) {
		client, cam, _ := newRoutedCamera(t, map[string]route{})

		if err := client.SetMotionSensitivity(context.Background(), cam, 50); err == nil || !strings.Contains(err.Error(), "404") {
			t.Errorf("err = %v", err)
		}
	})

	t.Run("sensitivity missing", func(t *testing.T) {
		client, cam, _ := newRoutedCamera(t, map[string]route{
			"GET " + path: {body: `<MotionDetection></MotionDetection>`},
		})

		if err := client.SetMotionSensitivity(context.Background(), cam, 50); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("put rejected", func(t *testing.T) {
		client, cam, _ := newRoutedCamera(t, map[string]route{
			"GET " + path: {body: current},
			"PUT " + path: {status: http.StatusForbidden, body: "forbidden"},
		})

		if err := client.SetMotionSensitivity(context.Background(), cam, 50); err == nil || !strings.Contains(err.Error(), "403") {
			t.Errorf("err = %v", err)
		}
	})
}

func TestSetThermometry(t *testing.T) {
	const path = "/ISAPI/Thermal/channels/1/thermometry/basicParam"
	const current = `<ThermometryBasicParam><id>1</id><enabled>false</enabled></ThermometryBasicParam>`

	t.Run("enable", func(t *testing.T) {
		client, cam, rc := newRoutedCamera(t, map[string]route{
			"GET " + path: {body: current},
			"PUT " + path: {body: responseStatusOK},
		})

		if err := client.SetThermometry(context.Background(), cam, true); err != nil {
			t.Fatal(err)
		}
		if body := rc.last(http.MethodPut).body; !strings.Contains(body, "<enabled>true</enabled>") {
			t.Errorf("PUT body = %s", body)
		}
	})

	t.Run("get fails", func(t *testing.T) {
		client, cam, _ := newRoutedCamera(t, map[string]route{})

		if err := client.SetThermometry(context.Background(), cam, true); err == nil || !strings.Contains(err.Error(), "404") {
			t.Errorf("err = %v", err)
		}
	})

	t.Run("enabled missing", func(t *testing.T) {
		client, cam, _ := newRoutedCamera(t, map[string]route{
			"GET " + path: {body: `<Other/>`},
		})

		if err := client.SetThermometry(context.Background(), cam, true); err == nil {
			t.Error("expected error")
		}
	})

	t.Run("put rejected", func(t *testing.T) {
		client, cam, _ := newRoutedCamera(t, map[string]route{
			"GET " + path: {body: current},
			"PUT " + path: {status: http.StatusForbidden, body: "forbidden"},
		})

		if err := client.SetThermometry(context.Background(), cam, true); err == nil || !strings.Contains(err.Error(), "403") {
			t.Errorf("err = %v", err)
		}
	})
}

func TestGetThermometryRulesTemperatureInfo(t *testing.T) {
	client, cam, rc := newRoutedCamera(t, map[string]route{
		"GET /ISAPI/Thermal/channels/1/thermometry/2/rulesTemperatureInfo": {
			body: `{"ThermometryRulesTemperatureInfoList":{"ThermometryRulesTemperatureInfo":[{"maxTemperature":40.5}]}}`,
		},
	})

	got, err := client.GetThermometryRulesTemperatureInfo(context.Background(), cam, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].MaxTemperature != 40.5 {
		t.Errorf("got %+v", got)
	}
	if q := rc.last(http.MethodGet).query; q != "format=json" {
		t.Errorf("query = %q, want format=json", q)
	}
}

func TestWaitStopsOnContextCancel(t *testing.T) {
	// The camera never reaches the requested position.
	ptz := &fakePTZ{azimuth: "0", elevation: "0", zoom: "1"}
	stuck := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			_, _ = io.Copy(io.Discard, r.Body)
			return
		}
		ptz.handler(w, r)
	}

	client, cam := newTestCamera(t, stuck)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	if _, err := client.SetZoom(ctx, cam, 20); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("SetZoom err = %v, want deadline exceeded", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel2()

	if _, err := client.SetAbsolute(ctx2, cam, model.PanTilt{X: 900, Y: 0}, 10); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("SetAbsolute err = %v, want deadline exceeded", err)
	}
}
