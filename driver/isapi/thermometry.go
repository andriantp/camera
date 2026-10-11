package isapi

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/andriantp/camera/driver/isapi/json"
	"github.com/andriantp/camera/model"
)

func (c *Client) EnableThermometry(ctx context.Context, camera Camera) error {
	return c.setThermometryEnabled(ctx, camera, true)
}

func (c *Client) DisableThermometry(ctx context.Context, camera Camera) error {
	return c.setThermometryEnabled(ctx, camera, false)
}

func (c *Client) setThermometryEnabled(ctx context.Context, camera Camera, enabled bool) error {
	path := fmt.Sprintf(
		"/ISAPI/Thermal/channels/%d/thermometry/basicParam",
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
		return fmt.Errorf("get thermometry basic param: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		return fmt.Errorf(
			"get thermometry basic param: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("read thermometry basic param: %w", err)
	}

	body, err = setThermometryEnabledXML(body, enabled)
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
		return fmt.Errorf("update thermometry basic param: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf(
			"update thermometry basic param: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	return nil
}

func setThermometryEnabledXML(body []byte, enabled bool) ([]byte, error) {
	re := regexp.MustCompile(
		`(?s)(<ThermometryBasicParam\b[^>]*>.*?<enabled>)[^<]*(</enabled>)`,
	)

	value := "false"
	if enabled {
		value = "true"
	}

	if !re.Match(body) {
		return nil, errors.New(
			"thermometry basic param: enabled element not found",
		)
	}

	return re.ReplaceAll(
		body,
		[]byte("${1}"+value+"${2}"),
	), nil
}

func (c *Client) GetThermometryRulesTemperatureInfo(ctx context.Context, camera Camera, rule int) ([]model.ThermometryRulesTemperatureInfo, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		thermometryRulesTemperatureInfoPath(camera.Channel, rule),
		nil,
	)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"get thermometry rules temperature info: status %d: %s",
			resp.StatusCode,
			body,
		)
	}

	return json.ParseThermometryRulesTemperatureInfo(resp)
}

func thermometryRulesTemperatureInfoPath(channel int, rule int) string {
	return fmt.Sprintf(
		"/ISAPI/Thermal/channels/%d/thermometry/%d/rulesTemperatureInfo?format=json",
		channel,
		rule,
	)
}
