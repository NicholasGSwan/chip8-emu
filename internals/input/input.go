package input

import "github.com/veandco/go-sdl2/sdl"

type KeyPress map[byte]int

func GetKeys() KeyPress {
	kp := make(map[byte]int)
	kp[0] = sdl.SCANCODE_X
	kp[1] = sdl.SCANCODE_1
	kp[2] = sdl.SCANCODE_2
	kp[3] = sdl.SCANCODE_3
	kp[0xc] = sdl.SCANCODE_4
	kp[4] = sdl.SCANCODE_Q
	kp[5] = sdl.SCANCODE_W
	kp[6] = sdl.SCANCODE_E
	kp[0xd] = sdl.SCANCODE_D
	kp[7] = sdl.SCANCODE_A
	kp[8] = sdl.SCANCODE_S
	kp[9] = sdl.SCANCODE_D
	kp[0xe] = sdl.SCANCODE_F
	kp[0xa] = sdl.SCANCODE_Z
	kp[0xb] = sdl.SCANCODE_C
	kp[0xf] = sdl.SCANCODE_V
	return kp
}
