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
	dev, err := c.device(camera)
	if err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}

	return c.getStatus(dev, camera.Profile)
}

func (c *Client) getStatus(dev *usego.Device, profile string) (*model.Status, error) {
	resp, err := dev.CallMethod(useptz.GetStatus{
		ProfileToken: onvif2.ReferenceToken(profile),
	})
	if err != nil {
		return nil, fmt.Errorf("get status: %w", err)
	}

	return xml.ParseStatus(resp)
}
