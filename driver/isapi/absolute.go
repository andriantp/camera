package isapi

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetAbsolute(ctx context.Context, camera Camera) (*model.Absolute, error) {
	return c.getAbsolute(ctx, camera)
}

func (c *Client) getAbsolute(ctx context.Context, camera Camera) (*model.Absolute, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		fmt.Sprintf(
			"/ISAPI/PTZCtrl/channels/%d/absoluteEx",
			camera.Channel,
		),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get absolute: %w", err)
	}

	return xml.ParseAbsolute(resp)
}

func (c *Client) SetAbsolute(ctx context.Context, camera Camera, position model.PanTilt, zoom float64) (*model.Absolute, error) {
	cap, err := c.getPTZCapabilities(ctx, camera)
	if err != nil {
		return nil, fmt.Errorf("read capabilities: %w", err)
	}

	if cap.AbsolutePanTilt == nil {
		return nil, fmt.Errorf("absolute pan tilt is not supported")
	}

	if position.X < cap.AbsolutePanTilt.XRange.Min ||
		position.X > cap.AbsolutePanTilt.XRange.Max {
		return nil, fmt.Errorf(
			"pan out of range min:%.2f max:%.2f",
			cap.AbsolutePanTilt.XRange.Min,
			cap.AbsolutePanTilt.XRange.Max,
		)
	}

	if position.Y < cap.AbsolutePanTilt.YRange.Min ||
		position.Y > cap.AbsolutePanTilt.YRange.Max {
		return nil, fmt.Errorf(
			"tilt out of range min:%.2f max:%.2f",
			cap.AbsolutePanTilt.YRange.Min,
			cap.AbsolutePanTilt.YRange.Max,
		)
	}

	if zoom < cap.Zoom.Absolute.Min ||
		zoom > cap.Zoom.Absolute.Max {
		return nil, fmt.Errorf(
			"zoom out of range min:%.2f max:%.2f",
			cap.Zoom.Absolute.Min,
			cap.Zoom.Absolute.Max,
		)
	}

	targetZoom := zoom
	requestZoom := zoom / 10

	targetPosition := position
	requestPosition := model.PanTilt{
		X: position.X / 10,
		Y: position.Y / 10,
	}

	body := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<PTZAbsoluteEx version="2.0" xmlns="http://www.hikvision.com/ver20/XMLSchema">
	<elevation>%.2f</elevation>
	<azimuth>%.2f</azimuth>
	<absoluteZoom>%.2f</absoluteZoom>
</PTZAbsoluteEx>`,
		requestPosition.Y,
		requestPosition.X,
		requestZoom,
	)
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
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, response)
	}

	return c.waitAbsolute(
		ctx,
		camera,
		targetPosition,
		targetZoom,
	)
}

func (c *Client) waitAbsolute(ctx context.Context, camera Camera, targetPosition model.PanTilt, targetZoom float64) (*model.Absolute, error) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	for {
		current, err := c.getAbsolute(ctx, camera)
		if err != nil {
			return nil, err
		}

		if almostEqual(current.Position.X, targetPosition.X) &&
			almostEqual(current.Position.Y, targetPosition.Y) &&
			almostEqual(current.Zoom, targetZoom) {
			return current, nil
		}

		select {
		case <-ticker.C:
		case <-timeout.C:
			return nil, fmt.Errorf(
				"timeout waiting absolute position: x=%.2f y=%.2f zoom=%.2f, current x=%.2f y=%.2f zoom=%.2f",
				targetPosition.X,
				targetPosition.Y,
				targetZoom,
				current.Position.X,
				current.Position.Y,
				current.Zoom,
			)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.1
}
