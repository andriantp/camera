package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetDeviceInfo(ctx context.Context, camera Camera) (*model.DeviceInfo, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		"/ISAPI/System/deviceInfo",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get device info: %w", err)
	}

	info, err := xml.ParseDeviceInfo(resp)
	if err != nil {
		return nil, fmt.Errorf("get device info: %w", err)
	}

	return info, nil
}

func (c *Client) Reboot(ctx context.Context, camera Camera) error {
	resp, err := c.request(
		ctx,
		http.MethodPut,
		camera,
		"/ISAPI/System/reboot",
		nil,
	)
	if err != nil {
		return fmt.Errorf("reboot camera: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"reboot camera: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	return nil
}
