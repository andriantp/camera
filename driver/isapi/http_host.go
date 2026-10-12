package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetHTTPHosts(ctx context.Context, camera Camera) ([]model.HTTPHostNotification, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		"/ISAPI/Event/notification/httpHosts",
		nil,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"get HTTP hosts: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	return xml.ParseHTTPHosts(resp)
}
