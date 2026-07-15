package display

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var window *sdl.Window
var rend *sdl.Renderer

const scale int32 = 16

func init() {
	var err error
	window, err = sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 64*scale, 32*scale, sdl.WINDOW_SHOWN)
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
	rend.SetDrawColor(0, 0, 0, 255)
	rend.Clear()
	rend.SetDrawColor(255, 255, 255, 255)
	points := make([]sdl.Point, 0)
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[0]); x++ {
			if pixels[y][x] {

				xDraw := x * int(scale)
				yDraw := y * int(scale)
				for i := 0; i < int(scale); i++ {
					for j := 0; j < int(scale); j++ {

						points = append(points, sdl.Point{X: int32(xDraw + i), Y: int32(yDraw + j)})
					}
				}
				//rend.DrawPoint(int32(x), int32(y))

			} else {
				// fmt.Print("00")
			}
		}
		// fmt.Print("\n")
	}
	rend.DrawPoints(points)
	rend.Present()
}

// brought this version back just to see if it rendered things any differently
func DrawDisplaypoints(pixels [32][64]bool) {

	//points := make([]sdl.Point, 0)
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[0]); x++ {
			if pixels[y][x] {
				rend.SetDrawColor(255, 255, 255, 255)
				rend.Clear()
				xDraw := x * int(scale)
				yDraw := y * int(scale)
				for i := 0; i < int(scale); i++ {
					for j := 0; j < int(scale); j++ {

						rend.DrawPoint(int32(xDraw+i), int32(yDraw+j))
					}
				}
				//rend.DrawPoint(int32(x), int32(y))

			} else {
				// fmt.Print("00")
			}
		}
		// fmt.Print("\n")
	}
	//rend.DrawPoints(points)
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
