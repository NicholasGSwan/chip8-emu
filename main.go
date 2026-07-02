package main

import (
	"fmt"

	emu "github.com/NicholasGSwan/chip8-emu/internals"
)

func main() {
	fmt.Println("Hello, this is the start of the chip 8 emu")
	emuMem := new(emu.EmuMemory)

	emuMem.RunEmu("IBM Logo.ch8")
	emu.DestroyResources()
}
