package rcon

import (
	"fmt"

	"github.com/paulsavides/palworld-manager/internal/clients"
)

func RunBroadcast(options clients.RconClientOptions) {
	client, err := clients.Rcon(options)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer client.Close()

	client.Broadcast("1 2 3 4 5 6 7 8 9 1 2 3 4 5 6 7 8 9 1 2 3 4 5 6 7 8 9 1 2 3 4 5 6 7 8 9")
}
