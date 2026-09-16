package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type deviceInfo struct {
	DeviceName        string `xml:"deviceName"`
	DeviceID          string `xml:"deviceID"`
	DeviceDescription string `xml:"deviceDescription"`
	DeviceLocation    string `xml:"deviceLocation"`
	Model             string `xml:"model"`
	SerialNumber      string `xml:"serialNumber"`
	FirmwareVersion   string `xml:"firmwareVersion"`
}

func ParseDeviceInfo(resp *http.Response) (*model.DeviceInfo, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", body)
	}

	var data deviceInfo
	if err := xml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	return &model.DeviceInfo{
		DeviceName:        data.DeviceName,
		DeviceID:          data.DeviceID,
		DeviceDescription: data.DeviceDescription,
		DeviceLocation:    data.DeviceLocation,
		Model:             data.Model,
		SerialNumber:      data.SerialNumber,
		FirmwareVersion:   data.FirmwareVersion,
	}, nil
}
