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

	// for op1 <= 0xf1 {
	// 	emuMem.Decode()
	// 	op1 += 0x10
	// }
	// emuMem.Memory[500] = 0x00
	// emuMem.Memory[501] = 0xe0
	// opcode := emuMem.Fetch()
	// emu.PrintOpCode(opcode)
	// // fmt.Printf("opcode values are : %x %x %x %x\n", opcode.opval, opcode[1], opcode[2], opcode[3])
	// emuMem.Decode(opcode)

	emuMem.LoadRom("IBM Logo.ch8")
	opCode := emuMem.Fetch()
	for opCode.GetThreeval() != 0 {
		emu.PrintOpCode(opCode)
		emuMem.Decode(opCode)
		emuMem.PrintPcCounter()
		opCode = emuMem.Fetch()
	}
	emu.DestroyResources()
}
