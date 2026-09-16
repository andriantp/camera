package onvif

import (
	"context"
	"fmt"

	usemedia "github.com/use-go/onvif/media"

	xml "github.com/andriantp/camera/driver/onvif/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetProfiles(ctx context.Context, camera Camera) ([]model.Profile, error) {
	dev, err := c.device(camera)
	if err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}

	resp, err := dev.CallMethod(usemedia.GetProfiles{})
	if err != nil {
		return nil, fmt.Errorf("get profiles: %w", err)
	}

	return xml.ParseProfiles(resp)
}
