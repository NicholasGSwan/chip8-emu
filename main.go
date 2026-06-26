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
	// emuMem.SetDelayTimer(200)
	// emuMem.Decrement()

	op1 := 0x01
	op2 := byte(0x00)
	for op1 <= 0xf1 {
		emuMem.Decode(byte(op1), op2)
		op1 += 0x10
	}
}
