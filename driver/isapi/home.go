package isapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) GotoHome(ctx context.Context, camera Camera) error {
	resp, err := c.request(
		ctx,
		http.MethodPut,
		camera,
		fmt.Sprintf(
			"/ISAPI/PTZCtrl/channels/%d/homeposition/goto",
			camera.Channel,
		),
		bytes.NewBufferString(`<PTZData></PTZData>`),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%d %s", resp.StatusCode, body)
	}

	return nil
}
