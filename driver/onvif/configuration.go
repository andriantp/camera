package onvif

import (
	"context"
	"fmt"

	useptz "github.com/use-go/onvif/ptz"

	"github.com/andriantp/camera/driver/onvif/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetConfigurations(ctx context.Context, camera Camera) (*model.Configuration, error) {
	dev, err := c.device(camera)
	if err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}

	resp, err := dev.CallMethod(useptz.GetConfigurations{})
	if err != nil {
		return nil, fmt.Errorf("get configurations: %w", err)
	}

	return xml.ParseConfiguration(resp)
}
