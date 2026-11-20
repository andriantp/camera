package onvif

import (
	"context"
	"fmt"

	usego "github.com/use-go/onvif"
	usemedia "github.com/use-go/onvif/media"

	"github.com/andriantp/camera/driver/onvif/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetProfiles(ctx context.Context, camera Camera) ([]model.Profile, error) {
	var result []model.Profile

	err := c.withDevice(ctx, camera, func(dev *usego.Device) error {
		resp, err := dev.CallMethod(usemedia.GetProfiles{})
		if err != nil {
			return fmt.Errorf("get profiles: %w", err)
		}

		result, err = xml.ParseProfiles(resp)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
