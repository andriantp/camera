package xml

import (
	"net/http"
	"reflect"
	"testing"
)

func TestParseSystemCapabilities(t *testing.T) {
	got, err := ParseSystemCapabilities(response(http.StatusOK, `<DeviceCap>
	<SysCap>
		<isSupportSubscribeEvent>true</isSupportSubscribeEvent>
		<NetworkCap><isSupportNtp>true</isSupportNtp><isSupportHttps>true</isSupportHttps></NetworkCap>
		<VideoCap><videoInputPortNums>1</videoInputPortNums><isSupportPrivacyMask>true</isSupportPrivacyMask></VideoCap>
		<AudioCap><audioInputNums>1</audioInputNums><audioOutputNums>2</audioOutputNums></AudioCap>
	</SysCap>
	<SmartCap><isSupportROI>true</isSupportROI><isSupportLineDetection>true</isSupportLineDetection></SmartCap>
	<voicetalkNums>1</voicetalkNums>
	<isSupportSnapshot>true</isSupportSnapshot>
</DeviceCap>`))
	if err != nil {
		t.Fatal(err)
	}

	if !got.Network.NTP || !got.Network.HTTPS || got.Network.FTP {
		t.Errorf("Network = %+v", got.Network)
	}
	if got.Video.InputPorts != 1 || !got.Video.PrivacyMask {
		t.Errorf("Video = %+v", got.Video)
	}
	if got.Audio.Input != 1 || got.Audio.Output != 2 {
		t.Errorf("Audio = %+v", got.Audio)
	}
	if !got.Snapshot.Supported || got.VoiceTalk.Channels != 1 || !got.Events.Subscribe {
		t.Errorf("Snapshot/VoiceTalk/Events = %+v %+v %+v", got.Snapshot, got.VoiceTalk, got.Events)
	}
	if !got.Analytics.ROI || !got.Analytics.LineDetection || got.Analytics.FaceDetection {
		t.Errorf("Analytics = %+v", got.Analytics)
	}
}

func TestParseStreamingCapabilities(t *testing.T) {
	got, err := ParseStreamingCapabilities(response(http.StatusOK, `<StreamingChannel>
	<Transport>
		<ControlProtocolList><ControlProtocol><streamingTransport>RTSP</streamingTransport></ControlProtocol></ControlProtocolList>
		<Unicast><enabled>true</enabled><rtpTransportType opt="RTP/UDP,RTP/TCP">RTP/TCP</rtpTransportType></Unicast>
		<Multicast><enabled>false</enabled></Multicast>
		<Security><enabled>true</enabled><certificateType opt="digest">digest</certificateType></Security>
	</Transport>
	<Video>
		<enabled>true</enabled>
		<videoInputChannelID>1</videoInputChannelID>
		<videoResolutionWidth>2560</videoResolutionWidth>
		<videoResolutionHeight>1440</videoResolutionHeight>
		<SVC><enabled>true</enabled></SVC>
	</Video>
	<Audio><enabled>true</enabled><audioInputChannelID>1</audioInputChannelID></Audio>
	<isSupportRTCPCfg>true</isSupportRTCPCfg>
</StreamingChannel>`))
	if err != nil {
		t.Fatal(err)
	}

	tr := got.Transport
	if !reflect.DeepEqual(tr.ControlProtocols, []string{"RTSP"}) || !tr.Unicast || tr.Multicast || !tr.Security {
		t.Errorf("Transport = %+v", tr)
	}
	if !reflect.DeepEqual(tr.RTPTransport.Options, []string{"RTP/UDP", "RTP/TCP"}) || tr.RTPTransport.Type != "RTP/TCP" {
		t.Errorf("RTPTransport = %+v", tr.RTPTransport)
	}
	if !reflect.DeepEqual(tr.CertificateTypes, []string{"digest"}) {
		t.Errorf("CertificateTypes = %v", tr.CertificateTypes)
	}
	if !got.Video.Enabled || !got.Video.SVC || !reflect.DeepEqual(got.Video.Resolutions, []string{"2560x1440"}) ||
		!reflect.DeepEqual(got.Video.InputChannels, []int{1}) {
		t.Errorf("Video = %+v", got.Video)
	}
	if !got.Audio.Enabled || !got.Features.RTCP {
		t.Errorf("Audio/Features = %+v %+v", got.Audio, got.Features)
	}
}

func TestParseStreaming(t *testing.T) {
	got, err := ParseStreaming(response(http.StatusOK, `<StreamingChannelList>
	<StreamingChannel>
		<id>101</id><channelName>Camera 01</channelName><enabled>true</enabled>
		<Video><videoCodecType>H.265</videoCodecType><videoResolutionWidth>2560</videoResolutionWidth><videoResolutionHeight>1440</videoResolutionHeight></Video>
		<Audio><enabled>true</enabled></Audio>
	</StreamingChannel>
	<StreamingChannel><id>102</id><Video><videoCodecType>H.264</videoCodecType></Video></StreamingChannel>
</StreamingChannelList>`))
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("channels = %d, want 2", len(got))
	}
	c := got[0]
	if c.ID != 101 || c.Name != "Camera 01" || c.Video.Codec != "H.265" ||
		c.Video.Resolution.Width != 2560 || c.Video.Resolution.Height != 1440 || !c.Audio.Enabled {
		t.Errorf("channel = %+v", c)
	}
	if got[1].ID != 102 || got[1].Audio.Enabled {
		t.Errorf("channel 2 = %+v", got[1])
	}
}

func TestParseImageCapabilities(t *testing.T) {
	got, err := ParseImageCapabilities(response(http.StatusOK, `<ImageChannel>
	<WDR><mode opt="close,open,auto">close</mode></WDR>
	<BLC><BLCMode opt="UP,DOWN">UP</BLCMode></BLC>
	<IrcutFilter><IrcutFilterType opt="auto,day,night">auto</IrcutFilterType></IrcutFilter>
	<WhiteBalance><WhiteBalanceStyle opt="auto,manual">auto</WhiteBalanceStyle></WhiteBalance>
	<Exposure><ExposureType opt="auto">auto</ExposureType></Exposure>
	<FocusConfiguration><focusStyle opt="AUTO,MANUAL">AUTO</focusStyle></FocusConfiguration>
	<EIS><enabled opt="true,false">false</enabled></EIS>
</ImageChannel>`))
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got.WDR.Modes, []string{"close", "open", "auto"}) {
		t.Errorf("WDR = %v", got.WDR.Modes)
	}
	if !reflect.DeepEqual(got.IRCutFilter.Types, []string{"auto", "day", "night"}) {
		t.Errorf("IRCutFilter = %v", got.IRCutFilter.Types)
	}
	if !reflect.DeepEqual(got.Focus.Styles, []string{"AUTO", "MANUAL"}) || !got.EIS.Supported {
		t.Errorf("Focus/EIS = %v %+v", got.Focus.Styles, got.EIS)
	}
	if got.Dehaze.Modes != nil {
		t.Errorf("missing Dehaze should give nil modes, got %v", got.Dehaze.Modes)
	}
}

func TestParseEventCapabilities(t *testing.T) {
	got, err := ParseEventCapabilities(response(http.StatusOK, `<EventCap>
	<isSupportHDFull>true</isSupportHDFull>
	<isSupportNicBroken>true</isSupportNicBroken>
	<isSupportIllAccess>true</isSupportIllAccess>
	<isSupportMotionDetection>true</isSupportMotionDetection>
	<isSupportAbnormalReboot>true</isSupportAbnormalReboot>
</EventCap>`))
	if err != nil {
		t.Fatal(err)
	}

	if !got.Storage.HDFull || got.Storage.HDError || !got.Network.NICBroken || !got.Access.IllegalAccess ||
		!got.Video.MotionDetection || got.Video.TamperDetection || !got.System.AbnormalReboot {
		t.Errorf("got %+v", got)
	}
}

func TestParseAnalyticsCapabilities(t *testing.T) {
	got, err := ParseAnalyticsCapabilities(response(http.StatusOK, `<SmartCap>
	<isSupportROI>true</isSupportROI>
	<isSupportFieldDetection>true</isSupportFieldDetection>
	<isSupportIntelliTrace>true</isSupportIntelliTrace>
</SmartCap>`))
	if err != nil {
		t.Fatal(err)
	}

	if !got.ROI || !got.FieldDetection || !got.IntelliTrace || got.FaceDetection {
		t.Errorf("got %+v", got)
	}
}

func TestParseContentCapabilities(t *testing.T) {
	got, err := ParseContentCapabilities(response(http.StatusOK, `<RacmCap>
	<nasNums>8</nasNums>
	<pictureSearchType opt="AllEvent,CMR,MOTION"/>
	<recordSearchType opt="time,event"/>
</RacmCap>`))
	if err != nil {
		t.Fatal(err)
	}

	if got.NASNums != 8 ||
		!reflect.DeepEqual(got.PictureSearchTypes, []string{"AllEvent", "CMR", "MOTION"}) ||
		!reflect.DeepEqual(got.RecordSearchTypes, []string{"time", "event"}) {
		t.Errorf("got %+v", got)
	}
}

func TestParseInputs(t *testing.T) {
	got, err := ParseInputs(response(http.StatusOK, `<VideoInputChannelList>
	<VideoInputChannel><id>1</id><inputPort>1</inputPort><name>Front Door</name><videoFormat>PAL</videoFormat></VideoInputChannel>
</VideoInputChannelList>`))
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || got[0].ID != 1 || got[0].InputPort != 1 || got[0].Name != "Front Door" || got[0].VideoFormat != "PAL" {
		t.Errorf("got %+v", got)
	}
}

func TestParsePresets(t *testing.T) {
	got, err := ParsePresets(response(http.StatusOK, `<PTZPresetList>
	<PTZPreset>
		<id>1</id><enabled>true</enabled><presetName>Gate</presetName>
		<AbsoluteHigh><elevation>100</elevation><azimuth>1800</azimuth><absoluteZoom>20</absoluteZoom></AbsoluteHigh>
	</PTZPreset>
	<PTZPreset><id>2</id><presetName>Parking</presetName></PTZPreset>
</PTZPresetList>`))
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("presets = %d, want 2", len(got))
	}
	p := got[0]
	if p.ID != 1 || p.Name != "Gate" || p.Position.Azimuth != 1800 || p.Position.Elevation != 100 || p.Position.Zoom != 20 {
		t.Errorf("preset = %+v", p)
	}
}

func TestParseDeviceInfo(t *testing.T) {
	got, err := ParseDeviceInfo(response(http.StatusOK, `<DeviceInfo>
	<deviceName>Gate Camera</deviceName>
	<deviceID>abc-123</deviceID>
	<model>DS-2DE4425IW</model>
	<serialNumber>SN0001</serialNumber>
	<firmwareVersion>V5.7.0</firmwareVersion>
</DeviceInfo>`))
	if err != nil {
		t.Fatal(err)
	}

	if got.DeviceName != "Gate Camera" || got.DeviceID != "abc-123" || got.Model != "DS-2DE4425IW" ||
		got.SerialNumber != "SN0001" || got.FirmwareVersion != "V5.7.0" {
		t.Errorf("got %+v", got)
	}
}

func TestParseStatus(t *testing.T) {
	got, err := ParseStatus(response(http.StatusOK, `<PTZStatus><AbsoluteHigh>
	<elevation>100</elevation><azimuth>1800</azimuth><absoluteZoom>20</absoluteZoom>
</AbsoluteHigh></PTZStatus>`))
	if err != nil {
		t.Fatal(err)
	}

	if got.Position.X != 1800 || got.Position.Y != 100 || got.Zoom != 20 {
		t.Errorf("got %+v", got)
	}
}

func TestParsersRejectInvalidXML(t *testing.T) {
	parsers := map[string]func(*http.Response) error{
		"ParseAbsolute":              func(r *http.Response) error { _, err := ParseAbsolute(r); return err },
		"ParseAnalyticsCapabilities": func(r *http.Response) error { _, err := ParseAnalyticsCapabilities(r); return err },
		"ParseContentCapabilities":   func(r *http.Response) error { _, err := ParseContentCapabilities(r); return err },
		"ParseContentSearch":         func(r *http.Response) error { _, err := ParseContentSearch(r); return err },
		"ParseDeviceInfo":            func(r *http.Response) error { _, err := ParseDeviceInfo(r); return err },
		"ParseEventCapabilities":     func(r *http.Response) error { _, err := ParseEventCapabilities(r); return err },
		"ParseHTTPHosts":             func(r *http.Response) error { _, err := ParseHTTPHosts(r); return err },
		"ParseImageCapabilities":     func(r *http.Response) error { _, err := ParseImageCapabilities(r); return err },
		"ParseInputs":                func(r *http.Response) error { _, err := ParseInputs(r); return err },
		"ParseMotionDetection":       func(r *http.Response) error { _, err := ParseMotionDetection(r); return err },
		"ParsePresets":               func(r *http.Response) error { _, err := ParsePresets(r); return err },
		"ParsePTZCapabilities":       func(r *http.Response) error { _, err := ParsePTZCapabilities(r); return err },
		"ParseStatus":                func(r *http.Response) error { _, err := ParseStatus(r); return err },
		"ParseStorage":               func(r *http.Response) error { _, err := ParseStorage(r); return err },
		"ParseStreaming":             func(r *http.Response) error { _, err := ParseStreaming(r); return err },
		"ParseStreamingCapabilities": func(r *http.Response) error { _, err := ParseStreamingCapabilities(r); return err },
		"ParseSystemCapabilities":    func(r *http.Response) error { _, err := ParseSystemCapabilities(r); return err },
	}

	for name, parse := range parsers {
		t.Run(name, func(t *testing.T) {
			if err := parse(response(http.StatusOK, "<broken")); err == nil {
				t.Error("expected error for invalid XML")
			}
		})
	}
}

func TestParsersErrorOnNon200(t *testing.T) {
	parsers := map[string]func(*http.Response) error{
		"ParseAnalyticsCapabilities": func(r *http.Response) error { _, err := ParseAnalyticsCapabilities(r); return err },
		"ParseContentCapabilities":   func(r *http.Response) error { _, err := ParseContentCapabilities(r); return err },
		"ParseEventCapabilities":     func(r *http.Response) error { _, err := ParseEventCapabilities(r); return err },
		"ParseImageCapabilities":     func(r *http.Response) error { _, err := ParseImageCapabilities(r); return err },
		"ParseInputs":                func(r *http.Response) error { _, err := ParseInputs(r); return err },
		"ParseStreamingCapabilities": func(r *http.Response) error { _, err := ParseStreamingCapabilities(r); return err },
	}

	for name, parse := range parsers {
		t.Run(name, func(t *testing.T) {
			resp, body := trackedResponse(http.StatusNotFound, "Not Found")
			if err := parse(resp); err == nil {
				t.Error("expected error for 404")
			}
			if !body.closed {
				t.Error("body was not closed")
			}
		})
	}
}
