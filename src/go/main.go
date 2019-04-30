package main

import (
	"fmt"

	"frank/src/go/controller"
)

func main() {
	controller.NewFrankController()

	var input string
	fmt.Scanln(&input)
	// fc.Start()
}
