package main

import (
	"fmt"

	"github.com/pagarme/warp-pipe/cmd"
)

func sample() string {
	return "sample"
}

func main() {

	if err := cmd.Execute(); err != nil {
		panic(err)
	}

	fmt.Println("Warp Pipe Project")
	fmt.Println(sample())
}
