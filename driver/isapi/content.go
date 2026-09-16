package isapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/andriantp/camera/driver/isapi/xml"
	"github.com/andriantp/camera/model"
)

func (c *Client) GetContentCapabilities(ctx context.Context, camera Camera) (*model.ContentCapabilities, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		"/ISAPI/ContentMgmt/capabilities",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get content capabilities: %w", err)
	}

	return xml.ParseContentCapabilities(resp)
}

func (c *Client) GetStorage(ctx context.Context, camera Camera) (*model.Storage, error) {
	resp, err := c.request(
		ctx,
		http.MethodGet,
		camera,
		"/ISAPI/ContentMgmt/Storage/hdd",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get storage: %w", err)
	}

	return xml.ParseStorage(resp)
}

type SearchRequest struct {
	SearchID             int       // identifier search request
	TrackID              int       // recording track
	StartTime            time.Time // Start Search time
	EndTime              time.Time // End Search time
	MaxResults           int       // total result
	SearchResultPosition int       // pagination/offset result
}

func (c *Client) ContentSearch(ctx context.Context, camera Camera, req SearchRequest) (*model.SearchResult, error) {
	xmlBody := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<CMSearchDescription>
	<searchID>%d</searchID>
	<trackID>%d</trackID>
	<timeSpanList>
		<timeSpan>
			<startTime>%s</startTime>
			<endTime>%s</endTime>
		</timeSpan>
	</timeSpanList>
	<maxResults>%d</maxResults>
	<searchResultPosition>%d</searchResultPosition>
</CMSearchDescription>`,
		req.SearchID,
		req.TrackID,
		req.StartTime.Format(time.RFC3339),
		req.EndTime.Format(time.RFC3339),
		req.MaxResults,
		req.SearchResultPosition,
	)

	resp, err := c.request(
		ctx,
		http.MethodPost,
		camera,
		"/ISAPI/ContentMgmt/search",
		strings.NewReader(xmlBody),
	)
	if err != nil {
		return nil, fmt.Errorf("get search: %w", err)
	}

	return xml.ParseContentSearch(resp)
}

func (c *Client) ContentDownload(ctx context.Context, camera Camera, playbackURI string) (int64, io.ReadCloser, error) {
	xmlPlaybackURI := strings.ReplaceAll(playbackURI, "&", "&amp;")

	xmlBody := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<downloadRequest version="1.0" xmlns="http://www.isapi.org/ver20/XMLSchema">
	<playbackURI>%s</playbackURI>
	<userName>%s</userName>
	<password>%s</password>
</downloadRequest>`,
		xmlPlaybackURI,
		camera.Username,
		camera.Password,
	)

	resp, err := c.requestDownload(
		ctx,
		http.MethodPost,
		camera,
		"/ISAPI/ContentMgmt/download",
		strings.NewReader(xmlBody),
	)
	if err != nil {
		return 0, nil, err
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		return 0, nil, fmt.Errorf("%d %s", resp.StatusCode, body)
	}

	return resp.ContentLength, resp.Body, nil
}
