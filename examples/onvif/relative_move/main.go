package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/andriantp/camera/examples/common"
	"github.com/andriantp/camera/model"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run . <LEFT|RIGHT|UP|DOWN>")
	}

	ctx := context.Background()
	client, cam := common.ONVIF()

	// status
	status, err := client.GetStatus(ctx, cam)
	if err != nil {
		log.Fatal(err)
	}
	js, err := json.MarshalIndent(status, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("prev status:%s\n", string(js))

	// profiles
	fmt.Printf("use profile:%s\n", cam.Profile) // got from profiles

	// move
	status, err = client.RelativeMove(ctx, cam, direction(strings.ToUpper(os.Args[1])), common.Wait)
	if err != nil {
		log.Fatal(err)
	}
	js, err = json.MarshalIndent(status, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("now status:%s\n", string(js))
}

func direction(turn string) model.PanTilt {
	log.Printf("move: %s", turn)
	switch turn {

	case "LEFT":
		return model.PanTilt{X: -0.05}

	case "RIGHT":
		return model.PanTilt{X: 0.05}

	case "UP":
		return model.PanTilt{Y: -0.05}

	case "DOWN":
		return model.PanTilt{Y: 0.05}

	default:
		log.Fatalf("invalid direction: %s", turn)
	}

	return model.PanTilt{}
}
