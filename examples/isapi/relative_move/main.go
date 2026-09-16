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
	client, cam := common.ISAPI()

	// move
	status, err := client.RelativeMove(ctx, cam, direction(strings.ToUpper(os.Args[1])), common.Wait)
	if err != nil {
		log.Fatal(err)
	}
	js, err := json.MarshalIndent(status, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("now status:%s\n", string(js))
}

func direction(turn string) model.Speed {
	log.Printf("move: %s", turn)
	switch turn {

	case "LEFT":
		return model.Speed{
			Pan:  -40,
			Tilt: 0,
		}

	case "RIGHT":
		return model.Speed{
			Pan:  40,
			Tilt: 0,
		}

	case "UP":
		return model.Speed{
			Pan:  0,
			Tilt: 40,
		}

	case "DOWN":
		return model.Speed{
			Pan:  0,
			Tilt: -40,
		}

	default:
		log.Fatalf("invalid direction: %s", turn)
	}

	return model.Speed{}
}
