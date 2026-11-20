package isapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/andriantp/camera/model"
)

// ErrMoveTimeout is returned (wrapped) when the camera does not reach the
// requested PTZ position within the configured camera.Config.MoveTimeout.
var ErrMoveTimeout = errors.New("timeout waiting for PTZ position")

const (
	defaultMoveTimeout = 5 * time.Second
	movePollInterval   = 100 * time.Millisecond
)

// waitPosition polls the camera's absolute position until reached reports
// true, the move timeout elapses, or ctx is done. target describes the
// requested position for the timeout error message.
func (c *Client) waitPosition(ctx context.Context, camera Camera, reached func(*model.Absolute) bool, target string) (*model.Absolute, error) {
	timeout := c.moveTimeout
	if timeout <= 0 {
		timeout = defaultMoveTimeout
	}

	ticker := time.NewTicker(movePollInterval)
	defer ticker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		current, err := c.getAbsolute(ctx, camera)
		if err != nil {
			return nil, err
		}

		if reached(current) {
			return current, nil
		}

		select {
		case <-ticker.C:
		case <-timer.C:
			return nil, fmt.Errorf(
				"%w after %v: target %s, current x=%.2f y=%.2f zoom=%.2f",
				ErrMoveTimeout,
				timeout,
				target,
				current.Position.X,
				current.Position.Y,
				current.Zoom,
			)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
