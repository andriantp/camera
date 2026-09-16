package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/andriantp/camera/examples/common"
)

func main() {
	client, cam := common.ONVIF()
	cap, err := client.GetCapabilities(context.Background(), cam)
	if err != nil {
		log.Fatal(err)
	}
	js, err := json.MarshalIndent(cap, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(js))
}
