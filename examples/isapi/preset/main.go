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
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run . <GET|GOTO>")
	}

	ctx := context.Background()

	client, cam := common.ISAPI()

	switch strings.ToUpper(os.Args[1]) {

	case "GET":
		body, err := client.GetPresets(ctx, cam)
		if err != nil {
			log.Fatal(err)
		}

		printJSON("GET PRESETS", body)

	case "GOTO":
		if len(os.Args) < 3 {
			log.Fatal("usage: go run . GOTO id")
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("invalid preset id")
		}

		fmt.Printf("========== GOTO %d ==========\n", id)
		if err := client.GotoPreset(ctx, cam, id); err != nil {
			log.Fatal(err)
		}
		fmt.Println("OK")

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
