package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/andriantp/camera/model"
)

func (c *Client) SetZoom(ctx context.Context, camera Camera, zoom float64) (*model.Absolute, error) {
	cap, err := c.getPTZCapabilities(ctx, camera)
	if err != nil {
		return nil, fmt.Errorf("read range: %w", err)
	}

	if cap.Zoom.Absolute.Min == 0 && cap.Zoom.Absolute.Max == 0 {
		return nil, fmt.Errorf("absolute zoom is not supported")
	}

	targetZoom := zoom
	requestZoom := zoom / 10

	if targetZoom < cap.Zoom.Absolute.Min || targetZoom > cap.Zoom.Absolute.Max {
		return nil, fmt.Errorf(
			"zoom out of range min:%.2f max:%.2f",
			cap.Zoom.Absolute.Min,
			cap.Zoom.Absolute.Max,
		)
	}

	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<PTZAbsoluteEx version="2.0" xmlns="http://www.hikvision.com/ver20/XMLSchema">
	<absoluteZoom>%.2f</absoluteZoom>
</PTZAbsoluteEx>`, requestZoom)

	resp, err := c.request(
		ctx,
		http.MethodPut,
		camera,
		fmt.Sprintf(
			"/ISAPI/PTZCtrl/channels/%d/absoluteEx",
			camera.Channel,
		),
		strings.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("set zoom: %w", err)
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, response)
	}

	return c.waitPosition(
		ctx,
		camera,
		func(current *model.Absolute) bool {
			return almostEqual(current.Zoom, targetZoom)
		},
		fmt.Sprintf("zoom=%.2f", targetZoom),
	)
}
