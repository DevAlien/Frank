package main

import (
	"frank/src/go/controller"
)

func main() {
	fc, _ := controller.NewFrankController()
	fc.Start()
}
