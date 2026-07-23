package main

import (
	"fmt"

	emu "github.com/NicholasGSwan/chip8-emu/internals"
)

func main() {
	fmt.Println("Hello, this is the start of the chip 8 emu")

	//emuMem.RunEmu("IBM Logo.ch8")
	// emuMem.RunEmu("test_opcode.ch8")
	// emuMem.RunEmu("test_opcode2.ch8")
	emu.RunEmu("TETRIS.ch8")
	//emuMem.RunEmu("BC_test.ch8")
	// emuMem.RunEmu("test.ch8")
	emu.DestroyResources()
}
