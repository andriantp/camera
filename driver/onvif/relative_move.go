package onvif

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	usego "github.com/use-go/onvif"
	useptz "github.com/use-go/onvif/ptz"
	onvif2 "github.com/use-go/onvif/xsd/onvif"

	"github.com/andriantp/camera/model"
)

func (c *Client) RelativeMove(ctx context.Context, camera Camera, panTilt model.PanTilt, wait time.Duration) (*model.Status, error) {
	err := c.withDevice(ctx, camera, func(dev *usego.Device) error {
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
			return fmt.Errorf("relative move: %w", err)
		}

		defer resp.Body.Close()

		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("relative move: status %d: %s", resp.StatusCode, bs)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Wait outside withDevice so other calls to this camera are not blocked.
	if wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
		}
	}

	return c.GetStatus(ctx, camera)
}
