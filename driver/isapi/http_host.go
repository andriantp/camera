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

func (c *Client) SetHTTPHost(ctx context.Context, camera Camera, host model.HTTPHostNotification, enabled bool) error {
	if host.ID <= 0 {
		return fmt.Errorf("HTTP host ID must be greater than 0")
	}

	if enabled && (host.IPAddress == "" || host.PortNo <= 0 || host.URL == "") {
		return fmt.Errorf("HTTP host URL, IP address, and port are required")
	}

	if !enabled {
		host.URL = "/"
		host.ProtocolType = "HTTP"
		host.ParameterFormatType = "XML"
		host.AddressingFormatType = "ipaddress"
		host.IPAddress = "0.0.0.0"
		host.PortNo = 80
		host.UserName = ""
		host.HTTPAuthenticationMethod = "none"
	}

	body, err := xml.MarshalHTTPHostNotification(host)
	if err != nil {
		return fmt.Errorf("marshal HTTP host notification: %w", err)
	}

	path := fmt.Sprintf(
		"/ISAPI/Event/notification/httpHosts/%d",
		host.ID,
	)

	resp, err := c.request(
		ctx,
		http.MethodPut,
		camera,
		path,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("set HTTP host: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"set HTTP host: status %d: %s",
			resp.StatusCode,
			responseBody,
		)
	}

	return xml.ParseResponseStatus(resp)
}
