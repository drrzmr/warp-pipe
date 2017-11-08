package main

import (
	"github.com/pagarme/warp-pipe/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
