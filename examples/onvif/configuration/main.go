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
	conf, err := client.GetConfigurations(context.Background(), cam)
	if err != nil {
		log.Fatal(err)
	}
	js, err := json.MarshalIndent(conf, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", string(js))
}
