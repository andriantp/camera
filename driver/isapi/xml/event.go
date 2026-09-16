package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type eventCapability struct {
	HDFull          bool `xml:"isSupportHDFull"`
	HDError         bool `xml:"isSupportHDError"`
	NICBroken       bool `xml:"isSupportNicBroken"`
	IPConflict      bool `xml:"isSupportIpConflict"`
	IllegalAccess   bool `xml:"isSupportIllAccess"`
	MotionDetection bool `xml:"isSupportMotionDetection"`
	TamperDetection bool `xml:"isSupportTamperDetection"`
	AbnormalReboot  bool `xml:"isSupportAbnormalReboot"`
}

func ParseEventCapabilities(resp *http.Response) (*model.EventCapability, error) {
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

	var data eventCapability

	if err := xml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal event capability: %w", err)
	}

	result := &model.EventCapability{
		Storage: model.StorageEventCapability{
			HDFull:  data.HDFull,
			HDError: data.HDError,
		},

		Network: model.NetworkEventCapability{
			NICBroken:  data.NICBroken,
			IPConflict: data.IPConflict,
		},

		Access: model.AccessEventCapability{
			IllegalAccess: data.IllegalAccess,
		},

		Video: model.VideoEventCapability{
			MotionDetection: data.MotionDetection,
			TamperDetection: data.TamperDetection,
		},

		System: model.SystemEventCapability{
			AbnormalReboot: data.AbnormalReboot,
		},
	}

	return result, nil
}
