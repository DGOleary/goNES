package main

import (
	"fmt"
	"goNES/nes"
)

func main() {
	bus := nes.CreateBus()

	fmt.Println(bus)
	fmt.Println(bus.CPU)
}
