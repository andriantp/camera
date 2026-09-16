package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andriantp/camera/examples/common"
)

func main() {
	ctx := context.Background()
	client, cam := common.ISAPI()

	fmt.Println("========== GOTO HOME ==========")
	if err := client.GotoHome(ctx, cam); err != nil {
		log.Fatal(err)
	}
	fmt.Println("OK")
}
