package onvif

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	useptz "github.com/use-go/onvif/ptz"
	onvif2 "github.com/use-go/onvif/xsd/onvif"

	"github.com/andriantp/camera/model"
)

func (c *Client) RelativeMove(ctx context.Context, camera Camera, panTilt model.PanTilt, wait time.Duration) (*model.Status, error) {
	dev, err := c.device(camera)
	if err != nil {
		return nil, fmt.Errorf("create device: %w", err)
	}

	resp, err := dev.CallMethod(useptz.RelativeMove{
		ProfileToken: onvif2.ReferenceToken(camera.Profile),
		Translation: onvif2.PTZVector{
			PanTilt: onvif2.Vector2D{
				X:     panTilt.X,
				Y:     panTilt.Y,
				Space: "http://www.onvif.org/ver10/tptz/PanTiltSpaces/TranslationGenericSpace",
			},
			Zoom: onvif2.Vector1D{
				X:     0,
				Space: "http://www.onvif.org/ver10/tptz/ZoomSpaces/TranslationGenericSpace",
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("relative move: %w", err)
	}

	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", bs)
	}

	if wait > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(wait):
		}
	}

	return c.getStatus(dev, camera.Profile)
}
