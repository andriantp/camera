package isapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/andriantp/camera/model"
)

// GetCapabilities collects all ISAPI capability groups. Not every camera
// supports every group, so a failing group is left nil and does not fail the
// call. An error is returned only when ctx is done or every group fails
// (e.g. camera unreachable or wrong credentials).
func (c *Client) GetCapabilities(ctx context.Context, camera Camera) (*model.Capabilities, error) {
	isapi := &model.IsapiCapabilities{}
	var errs []error

	collect := func(name string, err error) bool {
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
			return false
		}
		return true
	}

	if v, err := c.GetSystemCapabilities(ctx, camera); collect("system", err) {
		isapi.System = v
	}
	if v, err := c.GetStreamingCapabilities(ctx, camera); collect("streaming", err) {
		isapi.Streaming = v
	}
	if v, err := c.GetPTZCapabilities(ctx, camera); collect("ptz", err) {
		isapi.PTZ = v
	}
	if v, err := c.GetEventCapabilities(ctx, camera); collect("event", err) {
		isapi.Event = v
	}
	if v, err := c.GetImageCapabilities(ctx, camera); collect("image", err) {
		isapi.Image = v
	}
	if v, err := c.GetAnalyticsCapabilities(ctx, camera); collect("analytic", err) {
		isapi.Analytic = v
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const groups = 6
	if len(errs) == groups {
		return nil, fmt.Errorf("get capabilities: %w", errors.Join(errs...))
	}

	return &model.Capabilities{Isapi: isapi}, nil
}
