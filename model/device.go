package model

type DeviceInfo struct {
	DeviceName        string `json:"device_name"`
	DeviceID          string `json:"device_id"`
	DeviceDescription string `json:"device_description"`
	DeviceLocation    string `json:"device_location"`
	Model             string `json:"model"`
	SerialNumber      string `json:"serial_number"`
	FirmwareVersion   string `json:"firmware_version"`
}
