package isapi

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/andriantp/camera/model"
)

// GetCapabilities collects all ISAPI capability groups concurrently. Not every
// camera supports every group, so a failing group is left nil and does not
// fail the call. An error is returned only when ctx is done or every group
// fails (e.g. camera unreachable or wrong credentials).
func (c *Client) GetCapabilities(ctx context.Context, camera Camera) (*model.Capabilities, error) {
	isapi := &model.IsapiCapabilities{}

	// Each fetch writes only its own field of isapi.
	fetches := []struct {
		name string
		run  func() error
	}{
		{"system", func() (err error) { isapi.System, err = c.GetSystemCapabilities(ctx, camera); return }},
		{"streaming", func() (err error) { isapi.Streaming, err = c.GetStreamingCapabilities(ctx, camera); return }},
		{"ptz", func() (err error) { isapi.PTZ, err = c.GetPTZCapabilities(ctx, camera); return }},
		{"event", func() (err error) { isapi.Event, err = c.GetEventCapabilities(ctx, camera); return }},
		{"image", func() (err error) { isapi.Image, err = c.GetImageCapabilities(ctx, camera); return }},
		{"analytic", func() (err error) { isapi.Analytic, err = c.GetAnalyticsCapabilities(ctx, camera); return }},
	}

	errs := make([]error, len(fetches))

	var wg sync.WaitGroup
	for i, f := range fetches {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := f.run(); err != nil {
				errs[i] = fmt.Errorf("%s: %w", f.name, err)
			}
		}()
	}
	wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	failed := 0
	for _, err := range errs {
		if err != nil {
			failed++
		}
	}

	if failed == len(fetches) {
		return nil, fmt.Errorf("get capabilities: %w", errors.Join(errs...))
	}

	return &model.Capabilities{Isapi: isapi}, nil
}
