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
		log.Fatal("usage: go run . <INPUTS|STREAMING>")
	}

	ctx := context.Background()

	client, cam := common.ISAPI()

	switch strings.ToUpper(os.Args[1]) {

	case "INPUTS":
		body, err := client.GetInputs(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("INPUTS", body)

	case "STREAMING":
		body, err := client.GetStreaming(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}
		printJSON("STREAMING", body)

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
