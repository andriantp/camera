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

func (c *Client) GetPresets(ctx context.Context, camera Camera) ([]model.Preset, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		fmt.Sprintf("/ISAPI/PTZCtrl/channels/%d/presets", camera.Channel),
		nil,
	)
	if err != nil {
		return nil, err
	}

	return xml.ParsePresets(resp)
}

func (c *Client) GotoPreset(ctx context.Context, camera Camera, id int) error {

	body := fmt.Sprintf(
		`<PTZPreset version="2.0" xmlns="http://www.hikvision.com/ver20/XMLSchema">
    <id>%d</id>
</PTZPreset>`,
		id,
	)

	resp, err := c.request(
		ctx,
		http.MethodPut,
		camera,
		fmt.Sprintf(
			"/ISAPI/PTZCtrl/channels/%d/presets/%d/goto",
			camera.Channel,
			id,
		),
		bytes.NewBufferString(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s", resp.StatusCode, bs)
	}

	return nil
}
