package main

import (
	"fmt"
	"goNES/nes"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	bus := nes.CreateBus()
	cart := nes.CreateCartridge("../dk.nes")
	fmt.Println(cart)
	sdl.Init(sdl.INIT_VIDEO)
	window, renderer, err := sdl.CreateWindowAndRenderer(640, 480, 0)

	closeRequested := false

	//if the window or renderer couldn't initialize
	if err != nil {
		return
	}

	bus.PPU.ConnectRenderer(renderer)
	bus.InsertCartridge(cart)
	//bus.CPU.SetPC()
	bus.CPU.Reset()
	window.Show()

	renderer.Clear()
	// renderer.SetDrawColor(25, 100, 200, 255)
	// renderer.DrawPoint(0, 0)
	// renderer.DrawPoint(10, 90)
	// renderer.DrawPoint(100, 60)
	// renderer.DrawPoint(50, 30)
	// renderer.DrawPoint(200, 60)

	//TODO Code for making a surface, incase needed later
	// temp, _ := sdl.CreateRGBSurface(0, 50, 50, 32, 0, 0, 0, 0)
	// col := color.RGBA{R: 44, G: 213, B: 246, A: 255}
	// for i := 0; i < 50; i++ {
	// 	temp.Set(i, i, col)
	// }
	// tex, err := renderer.CreateTextureFromSurface(temp)
	// rect := sdl.Rect{X: 50, Y: 50, W: 50, H: 50}
	// renderer.Copy(tex, nil, &rect)
	// fmt.Println()
	for !closeRequested {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		renderer.Present()
		for !bus.PPU.Complete {
			bus.Clock()
		}
		bus.PPU.Complete = false
		time.Sleep(100 * time.Millisecond)

		renderer.Present()
		//make sure to put PollEvent to a variable because the rendering thread can go to nil mid-check and cause a null reference error
		event := sdl.PollEvent()
		if event != nil && event.GetType() == sdl.QUIT {
			closeRequested = true
		}
	}

	renderer.Destroy()
	window.Destroy()
	sdl.Quit()

}
