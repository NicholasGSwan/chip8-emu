package display

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var window *sdl.Window
var rend *sdl.Renderer

func init() {
	var err error
	window, err = sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 64, 32, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Printf("could not open window %s\n", err)
		panic(err)
	}

	rend, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Printf("coud not create renderer %s\n", err)
		panic(err)
	}

}

func DrawDisplay(pixels [32][64]bool) {
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[0]); x++ {
			if pixels[y][x] {
				rend.SetDrawColor(255, 255, 255, 255)
				rend.DrawPoint(int32(x), int32(y))

			} else {
				fmt.Print("00")
			}
		}
		fmt.Print("\n")
	}
	rend.Present()
}

func DestroyResources() {
	if window != nil {
		window.Destroy()
	}
	if rend != nil {
		rend.Destroy()
	}
}
