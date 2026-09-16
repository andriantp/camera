package xml

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/andriantp/camera/model"
)

type searchTypeXML struct {
	Types string `xml:"opt,attr"`
}

type contentCapabilitiesXML struct {
	NASNums           int           `xml:"nasNums"`
	PictureSearchType searchTypeXML `xml:"pictureSearchType"`
	RecordSearchType  searchTypeXML `xml:"recordSearchType"`
}

func ParseContentCapabilities(resp *http.Response) (*model.ContentCapabilities, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, body)
	}

	var src contentCapabilitiesXML

	if err := xml.Unmarshal(body, &src); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	return &model.ContentCapabilities{
		NASNums:            src.NASNums,
		PictureSearchTypes: splitOptions(src.PictureSearchType.Types),
		RecordSearchTypes:  splitOptions(src.RecordSearchType.Types),
	}, nil
}

func splitOptions(value string) []string {
	if value == "" {
		return nil
	}

	return strings.Split(value, ",")
}

type storageXML struct {
	HDD []hddXML `xml:"hdd"`
}

type hddXML struct {
	ID        int     `xml:"id"`
	Name      string  `xml:"hddName"`
	Status    string  `xml:"status"`
	Capacity  float64 `xml:"capacity"`
	FreeSpace float64 `xml:"freeSpace"`
}

func ParseStorage(resp *http.Response) (*model.Storage, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, body)
	}

	var src storageXML
	if err := xml.Unmarshal(body, &src); err != nil {
		return nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	if len(src.HDD) == 0 {
		return nil, fmt.Errorf("storage not found")
	}

	hdd := src.HDD[0]

	used := hdd.Capacity - hdd.FreeSpace

	usedPct := 0.0
	if hdd.Capacity > 0 {
		usedPct = math.Round(
			((hdd.Capacity-hdd.FreeSpace)/hdd.Capacity*100)*100,
		) / 100
	}

	return &model.Storage{
		Total:   hdd.Capacity,
		Used:    used,
		Free:    hdd.FreeSpace,
		Unit:    "MB",
		UsedPct: usedPct,
	}, nil
}

type cmSearchResultXML struct {
	SearchID          string `xml:"searchID"`
	ResponseStatus    bool   `xml:"responseStatus"`
	ResponseStatusStr string `xml:"responseStatusStrg"`
	TotalMatches      int    `xml:"totalMatches"`
	NumOfMatches      int    `xml:"numOfMatches"`

	MatchList struct {
		Items []searchMatchItemXML `xml:"searchMatchItem"`
	} `xml:"matchList"`
}

type searchMatchItemXML struct {
	TrackID  int `xml:"trackID"`
	TimeSpan struct {
		StartTime string `xml:"startTime"`
		EndTime   string `xml:"endTime"`
	} `xml:"timeSpan"`

	MediaSegmentDescriptor struct {
		ContentType string `xml:"contentType"`
		CodecType   string `xml:"codecType"`
		PlaybackURI string `xml:"playbackURI"`
	} `xml:"mediaSegmentDescriptor"`
}

func ParseContentSearch(resp *http.Response) (*model.SearchResult, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d %s", resp.StatusCode, body)
	}

	var src cmSearchResultXML

	if err := xml.Unmarshal(body, &src); err != nil {
		return nil, fmt.Errorf("unmarshal XML: %w", err)
	}

	matches := make([]model.Record, 0, len(src.MatchList.Items))

	for _, item := range src.MatchList.Items {
		startTime, err := time.Parse(time.RFC3339, item.TimeSpan.StartTime)
		if err != nil {
			return nil, fmt.Errorf("parse start time: %w", err)
		}

		endTime, err := time.Parse(time.RFC3339, item.TimeSpan.EndTime)
		if err != nil {
			return nil, fmt.Errorf("parse end time: %w", err)
		}

		sizeBytes, err := parsePlaybackURI(
			item.MediaSegmentDescriptor.PlaybackURI,
		)
		if err != nil {
			return nil, err
		}

		sizeMB := math.Round(float64(sizeBytes)/(1024*1024)*100) / 100

		matches = append(matches, model.Record{
			TrackID:  item.TrackID,
			Filename: buildFilename(startTime, endTime, ".mp4"),
			Timeframe: model.Timeframe{
				Start:    startTime,
				End:      endTime,
				Duration: formatDuration(endTime.Sub(startTime)),
			},
			Size: model.Size{
				Bytes: sizeBytes,
				MB:    sizeMB,
			},
			Type: model.Type{
				Content: item.MediaSegmentDescriptor.ContentType,
				Codec:   item.MediaSegmentDescriptor.CodecType,
			},
			PlaybackURI: &item.MediaSegmentDescriptor.PlaybackURI,
		})
	}

	return &model.SearchResult{
		SearchID:     src.SearchID,
		More:         src.ResponseStatusStr == "MORE",
		TotalMatches: src.TotalMatches,
		NumOfMatches: src.NumOfMatches,
		Matches:      matches,
	}, nil
}

func parsePlaybackURI(playbackURI string) (int64, error) {
	u, err := url.Parse(playbackURI)
	if err != nil {
		return 0, fmt.Errorf("parse playback URI: %w", err)
	}

	sizeStr := u.Query().Get("size")
	if sizeStr == "" {
		return 0, fmt.Errorf("size parameter not found")
	}

	sizeBytes, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse size: %w", err)
	}

	return sizeBytes, nil
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)

	hours := d / time.Hour
	d -= hours * time.Hour

	minutes := d / time.Minute
	d -= minutes * time.Minute

	seconds := d / time.Second

	var parts []string

	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}

	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}

	if seconds > 0 || len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}

func buildFilename(start, end time.Time, ext string) string {
	return fmt.Sprintf(
		"%s_%s%s",
		start.Format("20060102T150405Z"),
		end.Format("20060102T150405Z"),
		ext,
	)
}
