package main

import (
	"log"

	"github.com/oatmealraisin/tasker/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err.Error())
	}
}
