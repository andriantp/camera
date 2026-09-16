package isapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/andriantp/camera/model"
)

func (c *Client) RelativeMove(ctx context.Context, camera Camera, speed model.Speed, duration time.Duration) (*model.Status, error) {

	body := fmt.Sprintf(
		`<PTZData>
	<pan>%d</pan>
	<tilt>%d</tilt>
	<Momentary>
		<duration>%d</duration>
	</Momentary>
</PTZData>`,
		speed.Pan,
		speed.Tilt,
		duration.Milliseconds(),
	)

	resp, err := c.request(
		ctx,
		http.MethodPut,
		camera,
		fmt.Sprintf(
			"/ISAPI/PTZCtrl/channels/%d/momentary",
			camera.Channel,
		),
		bytes.NewBufferString(body),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, bs)
	}

	return c.getStatus(ctx, camera)
}
