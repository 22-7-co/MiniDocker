package main

import (
	"fmt"
	"mini-docker/Docker/cmd"
	"os"
)

func main() {

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
