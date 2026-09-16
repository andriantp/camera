package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andriantp/camera/examples/common"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run . <ALL|SYSTEM|STREAMING|PTZ|EVENT|IMAGE|ANALYTIC>")
	}

	ctx := context.Background()
	client, cam := common.ISAPI()

	switch strings.ToUpper(os.Args[1]) {

	case "DEVICE":
		body, err := client.GetDeviceInfo(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("GetDeviceInfo", body)

	case "ALL":
		body, err := client.GetCapabilities(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("ALL Capabilities", body)

	case "SYSTEM":
		body, err := client.GetSystemCapabilities(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("SYSTEM Capabilities", body)

	case "STREAMING":
		body, err := client.GetStreamingCapabilities(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("STREAMING Capabilities", body)

	case "PTZ":
		body, err := client.GetPTZCapabilities(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("PTZ Capabilities", body)

	case "EVENT":
		body, err := client.GetEventCapabilities(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("EVENT Capabilities", body)

	case "IMAGE":
		body, err := client.GetImageCapabilities(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("IMAGE Capabilities", body)

	case "ANALYTIC":
		body, err := client.GetAnalyticsCapabilities(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("ANALYTIC Capabilities", body)

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
