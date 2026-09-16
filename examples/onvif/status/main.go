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
	status, err := client.GetStatus(context.Background(), cam)
	if err != nil {
		log.Fatal(err)
	}
	js, err := json.MarshalIndent(status, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(js))
}
