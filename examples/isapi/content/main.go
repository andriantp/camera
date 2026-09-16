package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/andriantp/camera/driver/isapi"
	"github.com/andriantp/camera/examples/common"
)

func main() {
	ctx := context.Background()
	if len(os.Args) < 2 {
		log.Fatal("usage: go run . <CAPABILITIES|STORAGE|SEARCH|DOWNLOAD|SNAPSHOT>")
	}

	client, cam := common.ISAPI()

	switch strings.ToUpper(os.Args[1]) {
	case "CAPABILITIES":
		body, err := client.GetContentCapabilities(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("CAPABILITIES", body)

	case "STORAGE":
		body, err := client.GetStorage(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}
		printJSON("STORAGE", body)

	case "SEARCH":
		now := time.Now()

		startTime := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			0, 0, 0, 0,
			now.Location(),
		)

		endTime := startTime.Add(24*time.Hour - time.Second)

		req := isapi.SearchRequest{
			SearchID:             1,
			TrackID:              101,
			StartTime:            startTime,
			EndTime:              endTime,
			MaxResults:           2,
			SearchResultPosition: 0,
		}

		body, err := client.ContentSearch(ctx, cam, req)
		if err != nil {
			log.Fatal(err)
		}

		printJSONIter("SEARCH", body)

	case "DOWNLOAD":
		start := time.Date(2026, 8, 26, 16, 53, 51, 0, time.UTC)
		end := time.Date(2026, 8, 26, 17, 0, 1, 0, time.UTC)

		playbackURI := ""

		size, body, err := client.ContentDownload(ctx, cam, playbackURI)
		if err != nil {
			log.Fatal(err)
		}
		defer body.Close()

		filename := fmt.Sprintf(
			"%s-%s.mp4",
			start.Format("20060102T150405Z"),
			end.Format("20060102T150405Z"),
		)

		if err := downloadFile(body, filename, size); err != nil {
			log.Fatal(err)
		}

		fmt.Println("download completed successfully")

	case "SNAPSHOT":
		size, body, err := client.GetSnapshot(
			ctx,
			cam,
			isapi.SnapshotOptions{},
		)
		if err != nil {
			log.Fatal(err)
		}

		defer body.Close()

		filename := fmt.Sprintf("%d.jpg", time.Now().Unix())

		if err := downloadFile(body, filename, size); err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}

func printJSON(title string, v any) {
	js, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("========== %s ==========\n", title)
	fmt.Println(string(js))
}

func printJSONIter(title string, v any) {
	json := jsoniter.Config{
		EscapeHTML:             false,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()

	js, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("========== %s ==========\n", title)
	fmt.Println(string(js))
}

func downloadFile(body io.Reader, filename string, expectedSize int64) error {
	fmt.Printf("downloaded: %.2f MB\n", float64(expectedSize)/(1024*1024))

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	start := time.Now()
	var downloaded int64
	buf := make([]byte, 32*1024)
	for {
		n, err := body.Read(buf)

		if n > 0 {
			downloaded += int64(n)

			if _, werr := file.Write(buf[:n]); werr != nil {
				return fmt.Errorf("write file: %w", werr)
			}

			if expectedSize > 0 {
				fmt.Printf(
					"\rDownloaded %d / %d bytes (%.2f%%)",
					downloaded,
					expectedSize,
					float64(downloaded)*100/float64(expectedSize),
				)
			} else {
				fmt.Printf(
					"\rDownloaded %d bytes",
					downloaded,
				)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("read body: %w", err)
		}
	}

	fmt.Println()
	duration := time.Since(start)
	speedMB := float64(downloaded) / (1024 * 1024) / duration.Seconds()

	fmt.Printf("downloaded: %d bytes\n", downloaded)
	fmt.Printf("expected:   %d bytes\n", expectedSize)
	fmt.Printf("duration:   %s\n", duration.Round(time.Millisecond))
	fmt.Printf("speed:    %.2f MB/s\n", speedMB)

	if expectedSize > 0 && downloaded != expectedSize {
		return fmt.Errorf(
			"invalid download size: expected %d, got %d",
			expectedSize,
			downloaded,
		)
	}

	return nil
}
