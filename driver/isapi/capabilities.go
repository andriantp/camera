package isapi

import (
	"context"

	"github.com/andriantp/camera/model"
)

func (c *Client) GetCapabilities(ctx context.Context, camera Camera) (*model.Capabilities, error) {
	var capabilities model.Capabilities

	capabilities.Isapi = &model.IsapiCapabilities{}

	system, err := c.GetSystemCapabilities(ctx, camera)
	if err == nil {
		capabilities.Isapi.System = system
	}

	streaming, err := c.GetStreamingCapabilities(ctx, camera)
	if err == nil {
		capabilities.Isapi.Streaming = streaming
	}

	ptz, err := c.GetPTZCapabilities(ctx, camera)
	if err == nil {
		capabilities.Isapi.PTZ = ptz
	}

	event, err := c.GetEventCapabilities(ctx, camera)
	if err == nil {
		capabilities.Isapi.Event = event
	}

	image, err := c.GetImageCapabilities(ctx, camera)
	if err == nil {
		capabilities.Isapi.Image = image
	}

	analytic, err := c.GetAnalyticsCapabilities(ctx, camera)
	if err == nil {
		capabilities.Isapi.Analytic = analytic
	}

	return &capabilities, nil
}
