package isapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetMotionDetection(ctx context.Context, camera Camera) (*model.MotionDetection, error) {
	path := fmt.Sprintf(
		"/ISAPI/System/Video/inputs/channels/%d/motionDetection",
		camera.Channel,
	)

	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		path,
		nil,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"get motion detection: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	return xml.ParseMotionDetection(resp)
}

func (c *Client) SetVMDSensitivity(ctx context.Context, camera Camera, level int) error {
	if level < 0 || level > 100 {
		return fmt.Errorf("VMD sensitivity must be between 0 and 100")
	}

	path := fmt.Sprintf(
		"/ISAPI/System/Video/inputs/channels/%d/motionDetection",
		camera.Channel,
	)

	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		path,
		nil,
	)
	if err != nil {
		return fmt.Errorf("get motion detection: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		return fmt.Errorf(
			"get motion detection: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("read motion detection: %w", err)
	}

	body, err = xml.SetVMDSensitivity(body, level)
	if err != nil {
		return err
	}

	resp, err = c.request(
		ctx,
		http.MethodPut,
		camera,
		path,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("set VMD sensitivity: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"set VMD sensitivity: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	return xml.ParseResponseStatus(resp)
}
