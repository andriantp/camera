package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/andriantp/camera/examples/common"
	"github.com/andriantp/camera/model"
)

func main() {
	ctx := context.Background()
	if len(os.Args) < 2 {
		log.Fatal("usage: go run . <RANGE|CURRENT|GETABSOLUTE|ZOOM>")
	}

	client, cam := common.ISAPI()

	switch strings.ToUpper(os.Args[1]) {
	case "RANGE":
		body, err := client.GetPTZCapabilities(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("Range Zoom", body.Zoom)

	case "CURRENT":
		body, err := client.GetStatus(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("Current Absolute:", body)

	case "GETABSOLUTE":
		body, err := client.GetAbsolute(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("GET ABSOLUTE:", body)

	case "ZOOM":
		if len(os.Args) < 3 {
			log.Fatal("usage: go run . <ZOOM> <VALUE>")
		}

		value := os.Args[2]
		value = strings.ReplaceAll(value, " ", "")
		set, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal(err)
		}
		body, err := client.SetZoom(ctx, cam, float64(set))
		if err != nil {
			log.Fatal(err)
		}
		printJSON("ZOOM:", body)

	case "SETABSOLUTE":
		position := model.PanTilt{
			X: 1800,
			Y: 450,
		}

		absolute, err := client.SetAbsolute(
			ctx,
			cam,
			position,
			200,
		)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("SETABSOLUTE:", absolute)

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
