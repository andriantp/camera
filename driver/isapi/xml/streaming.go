package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/andriantp/camera/model"
)

type streamingChannelListXML struct {
	Channels []streamingChannelXML `xml:"StreamingChannel"`
}

type streamingChannelXML struct {
	ID      int    `xml:"id"`
	Name    string `xml:"channelName"`
	Enabled bool   `xml:"enabled"`

	Video videoXML `xml:"Video"`
	Audio audioXML `xml:"Audio"`
}

type videoXML struct {
	Enabled bool `xml:"enabled"`

	Codec string `xml:"videoCodecType"`

	ResolutionWidth  int `xml:"videoResolutionWidth"`
	ResolutionHeight int `xml:"videoResolutionHeight"`
}

type audioXML struct {
	Enabled bool `xml:"enabled"`
}

func ParseStreaming(resp *http.Response) ([]model.StreamingChannel, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}

	var src streamingChannelListXML
	if err := xml.Unmarshal(body, &src); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	channels := make([]model.StreamingChannel, 0, len(src.Channels))

	for _, ch := range src.Channels {
		channels = append(channels, model.StreamingChannel{
			ID:   ch.ID,
			Name: ch.Name,
			Video: model.Video{
				Codec: ch.Video.Codec,
				Resolution: model.Resolution{
					Width:  ch.Video.ResolutionWidth,
					Height: ch.Video.ResolutionHeight,
				},
			},
			Audio: model.Audio{
				Enabled: ch.Audio.Enabled,
			},
		})
	}

	return channels, nil
}

type streamingCapabilities struct {
	Transport streamingTransport `xml:"Transport"`
	Video     streamingVideo     `xml:"Video"`
	Audio     streamingAudio     `xml:"Audio"`

	DynamicCapability             bool `xml:"isSpportDynamicCapWithCondition"`
	SmartCodeWithoutRestart       bool `xml:"isSupportSmartCodeWithoutReStart"`
	RTCP                          bool `xml:"isSupportRTCPCfg"`
	SpecifiedImageQualitySnapshot bool `xml:"isSupportSpecifiedImageQualitySnap"`
}

type streamingTransport struct {
	ControlProtocolList streamingControlProtocolList `xml:"ControlProtocolList"`
	Multicast           streamingMulticast           `xml:"Multicast"`
	Unicast             streamingUnicast             `xml:"Unicast"`
	Security            streamingSecurity            `xml:"Security"`
}

type streamingControlProtocolList struct {
	ControlProtocol streamingControlProtocol `xml:"ControlProtocol"`
}

type streamingControlProtocol struct {
	Transport string `xml:"streamingTransport"`
	Opt       string `xml:"streamingTransport,attr"`
}

type streamingMulticast struct {
	Enabled bool `xml:"enabled"`
}

type streamingUnicast struct {
	Enabled          bool                  `xml:"enabled"`
	RTPTransportType streamingRTPTransport `xml:"rtpTransportType"`
}

type streamingRTPTransport struct {
	Options string `xml:"opt,attr"`
	Type    string `xml:",chardata"`
}

type streamingSecurity struct {
	Enabled         bool                     `xml:"enabled"`
	CertificateType streamingCertificateType `xml:"certificateType"`
}

type streamingCertificateType struct {
	Value string `xml:",chardata"`
	Opt   string `xml:"opt,attr"`
}

type streamingVideo struct {
	Enabled        bool `xml:"enabled"`
	InputChannelID int  `xml:"videoInputChannelID"`

	CodecType    string `xml:"videoCodecType"`
	CodecTypeOpt string `xml:"videoCodecType,attr"`

	ScanType    string `xml:"videoScanType"`
	ScanTypeOpt string `xml:"videoScanType,attr"`

	ResolutionWidth    int    `xml:"videoResolutionWidth"`
	ResolutionWidthOpt string `xml:"videoResolutionWidth,attr"`

	ResolutionHeight    int    `xml:"videoResolutionHeight"`
	ResolutionHeightOpt string `xml:"videoResolutionHeight,attr"`

	QualityControlType    string `xml:"videoQualityControlType"`
	QualityControlTypeOpt string `xml:"videoQualityControlType,attr"`

	MaxFrameRate    int    `xml:"maxFrameRate"`
	MaxFrameRateOpt string `xml:"maxFrameRate,attr"`

	SnapshotImageType    string `xml:"snapShotImageType"`
	SnapshotImageTypeOpt string `xml:"snapShotImageType,attr"`

	H264Profile    string `xml:"H264Profile"`
	H264ProfileOpt string `xml:"H264Profile,attr"`

	H265Profile    string `xml:"H265Profile"`
	H265ProfileOpt string `xml:"H265Profile,attr"`

	SVC streamingSVC `xml:"SVC"`
}

type streamingSVC struct {
	Enabled bool   `xml:"enabled"`
	Mode    string `xml:"SVCMode"`
}

type streamingAudio struct {
	Enabled        bool `xml:"enabled"`
	InputChannelID int  `xml:"audioInputChannelID"`

	CompressionType    string `xml:"audioCompressionType"`
	CompressionTypeOpt string `xml:"audioCompressionType,attr"`
}

func ParseStreamingCapabilities(resp *http.Response) (*model.StreamingCapabilities, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, body)
	}

	var data streamingCapabilities
	if err := xml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal streaming capabilities: %w", err)
	}

	result := &model.StreamingCapabilities{
		Transport: model.TransportCapabilities{
			ControlProtocols: []string{
				data.Transport.ControlProtocolList.ControlProtocol.Transport,
			},
			Multicast: data.Transport.Multicast.Enabled,
			Unicast:   data.Transport.Unicast.Enabled,
			RTPTransport: model.RTPTransport{
				Options: splitOptions(data.Transport.Unicast.RTPTransportType.Options),
				Type:    data.Transport.Unicast.RTPTransportType.Type,
			},
			Security: data.Transport.Security.Enabled,
			CertificateTypes: []string{
				data.Transport.Security.CertificateType.Value,
			},
		},

		Video: model.StreamingVideoCapabilities{
			Enabled: data.Video.Enabled,

			InputChannels: []int{
				data.Video.InputChannelID,
			},

			Resolutions: []string{
				fmt.Sprintf(
					"%dx%d",
					data.Video.ResolutionWidth,
					data.Video.ResolutionHeight,
				),
			},

			SVC: data.Video.SVC.Enabled,
		},

		Audio: model.StreamingAudioCapabilities{
			Enabled: data.Audio.Enabled,

			InputChannels: []int{
				data.Audio.InputChannelID,
			},
		},

		Features: model.StreamingFeatures{
			DynamicCapability:             data.DynamicCapability,
			SmartCodeWithoutRestart:       data.SmartCodeWithoutRestart,
			RTCP:                          data.RTCP,
			SpecifiedImageQualitySnapshot: data.SpecifiedImageQualitySnapshot,
		},
	}

	return result, nil
}

func parseIntList(value string) []int {
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]int, 0, len(parts))

	for _, part := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			continue
		}

		result = append(result, n)
	}

	return result
}
