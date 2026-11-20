package onvif

import (
	"context"
	"fmt"

	usego "github.com/use-go/onvif"
	useptz "github.com/use-go/onvif/ptz"
	onvif2 "github.com/use-go/onvif/xsd/onvif"

	"github.com/andriantp/camera/driver/onvif/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetStatus(ctx context.Context, camera Camera) (*model.Status, error) {
	var result *model.Status

	err := c.withDevice(ctx, camera, func(dev *usego.Device) error {
		var err error
		result, err = getStatus(dev, camera.Profile)
		return err
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getStatus(dev *usego.Device, profile string) (*model.Status, error) {
	resp, err := dev.CallMethod(useptz.GetStatus{
		ProfileToken: onvif2.ReferenceToken(profile),
	})
	if err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}

	return xml.ParseStatus(resp)
}
