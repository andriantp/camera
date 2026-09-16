package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/andriantp/camera/model"
)

type imageCapabilities struct {
	WDR          imageWDR          `xml:"WDR"`
	BLC          imageBLC          `xml:"BLC"`
	IRCutFilter  imageIRCutFilter  `xml:"IrcutFilter"`
	WhiteBalance imageWhiteBalance `xml:"WhiteBalance"`
	Exposure     imageExposure     `xml:"Exposure"`
	Focus        imageFocus        `xml:"FocusConfiguration"`
	EIS          imageEIS          `xml:"EIS"`
	Dehaze       imageDehaze       `xml:"Dehaze"`
}

type imageWDR struct {
	Mode imageOption `xml:"mode"`
}

type imageBLC struct {
	Mode imageOption `xml:"BLCMode"`
}

type imageIRCutFilter struct {
	Type imageOption `xml:"IrcutFilterType"`
}

type imageWhiteBalance struct {
	Style imageOption `xml:"WhiteBalanceStyle"`
}

type imageExposure struct {
	Type imageOption `xml:"ExposureType"`
}

type imageFocus struct {
	Style imageOption `xml:"focusStyle"`
}

type imageEIS struct {
	Enabled imageOption `xml:"enabled"`
}

type imageDehaze struct {
	Mode imageOption `xml:"DehazeMode"`
}

type imageOption struct {
	Value string `xml:",chardata"`
	Opt   string `xml:"opt,attr"`
}

func ParseImageCapabilities(resp *http.Response) (*model.ImageCapabilities, error) {
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

	var data imageCapabilities

	if err := xml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal image capabilities: %w", err)
	}

	return &model.ImageCapabilities{
		WDR: model.WDRCapabilities{
			Modes: splitOptions(data.WDR.Mode.Opt),
		},

		BLC: model.BLCCapabilities{
			Modes: splitOptions(data.BLC.Mode.Opt),
		},

		IRCutFilter: model.IRCutFilterCapabilities{
			Types: splitOptions(data.IRCutFilter.Type.Opt),
		},

		WhiteBalance: model.WhiteBalanceCapabilities{
			Styles: splitOptions(data.WhiteBalance.Style.Opt),
		},

		Exposure: model.ExposureCapabilities{
			Types: splitOptions(data.Exposure.Type.Opt),
		},

		Focus: model.FocusCapabilities{
			Styles: splitOptions(data.Focus.Style.Opt),
		},

		EIS: model.EISCapabilities{
			Supported: strings.Contains(data.EIS.Enabled.Opt, "true"),
		},

		Dehaze: model.DehazeCapabilities{
			Modes: splitOptions(data.Dehaze.Mode.Opt),
		},
	}, nil
}
