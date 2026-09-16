package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type systemCapabilities struct {
	SysCap    systemSysCap   `xml:"SysCap"`
	SmartCap  systemSmartCap `xml:"SmartCap"`
	VoiceTalk int            `xml:"voicetalkNums"`
	Snapshot  bool           `xml:"isSupportSnapshot"`
}

type systemSysCap struct {
	NetworkCap systemNetworkCap `xml:"NetworkCap"`
	VideoCap   systemVideoCap   `xml:"VideoCap"`
	AudioCap   systemAudioCap   `xml:"AudioCap"`

	SubscribeEvent bool `xml:"isSupportSubscribeEvent"`
}

type systemNetworkCap struct {
	NTP   bool `xml:"isSupportNtp"`
	FTP   bool `xml:"isSupportFtp"`
	UPnP  bool `xml:"isSupportUpnp"`
	DDNS  bool `xml:"isSupportDdns"`
	HTTPS bool `xml:"isSupportHttps"`
}

type systemVideoCap struct {
	InputPorts  int `xml:"videoInputPortNums"`
	OutputPorts int `xml:"videoOutputPortNums"`

	Heatmap     bool `xml:"isSupportHeatmap"`
	Counting    bool `xml:"isSupportCounting"`
	Picture     bool `xml:"isSupportPicture"`
	PrivacyMask bool `xml:"isSupportPrivacyMask"`
}

type systemAudioCap struct {
	Input  int `xml:"audioInputNums"`
	Output int `xml:"audioOutputNums"`
}

type systemSmartCap struct {
	ROI            bool `xml:"isSupportROI"`
	AudioDetection bool `xml:"isSupportAudioDetection"`
	FaceDetection  bool `xml:"isSupportFaceDetect"`
	LineDetection  bool `xml:"isSupportLineDetection"`
	FieldDetection bool `xml:"isSupportFieldDetection"`
	RegionEntrance bool `xml:"isSupportRegionEntrance"`
	RegionExiting  bool `xml:"isSupportRegionExiting"`
}

func ParseSystemCapabilities(resp *http.Response) (*model.SystemCapabilities, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}

	var data systemCapabilities
	if err := xml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	result := &model.SystemCapabilities{
		Network: model.NetworkCapabilities{
			NTP:   data.SysCap.NetworkCap.NTP,
			FTP:   data.SysCap.NetworkCap.FTP,
			UPnP:  data.SysCap.NetworkCap.UPnP,
			DDNS:  data.SysCap.NetworkCap.DDNS,
			HTTPS: data.SysCap.NetworkCap.HTTPS,
		},

		Video: model.VideoCapabilities{
			InputPorts:  data.SysCap.VideoCap.InputPorts,
			OutputPorts: data.SysCap.VideoCap.OutputPorts,
			Heatmap:     data.SysCap.VideoCap.Heatmap,
			Counting:    data.SysCap.VideoCap.Counting,
			Picture:     data.SysCap.VideoCap.Picture,
			PrivacyMask: data.SysCap.VideoCap.PrivacyMask,
		},

		Audio: model.AudioCapabilities{
			Input:  data.SysCap.AudioCap.Input,
			Output: data.SysCap.AudioCap.Output,
		},

		Snapshot: model.SnapshotCapabilities{
			Supported: data.Snapshot,
		},

		VoiceTalk: model.VoiceTalkCapabilities{
			Channels: data.VoiceTalk,
		},

		Events: model.EventCapabilities{
			Subscribe: data.SysCap.SubscribeEvent,
		},

		Analytics: model.AnalyticsCapabilities{
			ROI:            data.SmartCap.ROI,
			AudioDetection: data.SmartCap.AudioDetection,
			FaceDetection:  data.SmartCap.FaceDetection,
			LineDetection:  data.SmartCap.LineDetection,
			FieldDetection: data.SmartCap.FieldDetection,
			RegionEntrance: data.SmartCap.RegionEntrance,
			RegionExiting:  data.SmartCap.RegionExiting,
		},
	}

	return result, nil
}
