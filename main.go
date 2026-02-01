package main

import (
	"fmt"

	emu "github.com/NicholasGSwan/chip8-emu/internals"
)

func main() {
	fmt.Println("Hello, this is the start of the chip 8 emu")
	emuMem := new(emu.EmuMemory)
	emuMem.Init(emu.Font)

	for i := range len(emu.Font) {
		fmt.Printf("The val at %d is %v \n", i, emuMem.Memory[i])
	}
	emuMem.SetDelayTimer(200)
	emuMem.Decrement()
}
