package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/andriantp/camera/examples/common"
)

func main() {
	ctx := context.Background()

	client, cam := common.ISAPI()
	body, err := client.GetStatus(ctx, cam)
	if err != nil {
		log.Fatal(err)
	}

	printJSON("Status", body)
}

func printJSON(title string, v any) {
	js, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("========== %s ==========\n", title)
	fmt.Println(string(js))
}
