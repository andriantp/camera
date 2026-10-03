package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/driver/isapi/json"
	"github.com/andriantp/camera/model"
)

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
