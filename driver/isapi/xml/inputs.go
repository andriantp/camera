package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	"github.com/andriantp/camera/model"
)

type videoInputChannelListXML struct {
	Channels []videoInputChannelXML `xml:"VideoInputChannel"`
}

type videoInputChannelXML struct {
	ID          int    `xml:"id"`
	InputPort   int    `xml:"inputPort"`
	Name        string `xml:"name"`
	VideoFormat string `xml:"videoFormat"`
}

func ParseInputs(resp *http.Response) ([]model.Input, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, body)
	}

	var src videoInputChannelListXML
	if err := xml.Unmarshal(body, &src); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	inputs := make([]model.Input, 0, len(src.Channels))

	for _, ch := range src.Channels {
		inputs = append(inputs, model.Input{
			ID:          ch.ID,
			InputPort:   ch.InputPort,
			Name:        ch.Name,
			VideoFormat: ch.VideoFormat,
		})
	}

	return inputs, nil
}
