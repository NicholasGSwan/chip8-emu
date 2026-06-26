package emu

import (
	"fmt"
	"time"
)

type EmuMemory struct {
	Memory         [4096]byte
	programCounter int
	indexRegister  int
	stack          []int
	delayTimer     uint8
	soundTimer     uint8
	varRegister    [16]byte
}

var (
	Font = []byte{0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80}
)

func (mem *EmuMemory) Init(font []byte) {
	for i, b := range font {
		mem.Memory[i] = b
	}
}

func (mem *EmuMemory) Fetch() (byte, byte) {
	op1 := mem.Memory[mem.programCounter]
	op2 := mem.Memory[mem.programCounter+1]
	mem.programCounter = mem.programCounter + 2
	return op1, op2
}

func (mem *EmuMemory) incrementPc() {
	mem.programCounter += 2
}

func (mem *EmuMemory) Decode(op1, op2 byte) {
	nib1 := op1 & 0xF0
	//nib2 := op1 & 0x0F
	switch nib1 {
	case 0x00:
		fmt.Println("0")
	case 0x10:
		fmt.Println("1")
	case 0x20:
		fmt.Println("2")
	case 0x30:
		fmt.Println("3")
	case 0x40:
		fmt.Println("4")
	case 0x50:
		fmt.Println("5")
	case 0x60:
		fmt.Println("6")
	case 0x70:
		fmt.Println("7")
	case 0x80:
		fmt.Println("8")
	case 0x90:
		fmt.Println("9")
	case 0xA0:
		fmt.Println("10")
	case 0xB0:
		fmt.Println("11")
	case 0xc0:
		fmt.Println("12")
	case 0xd0:
		fmt.Println("13")
	case 0xe0:
		fmt.Println("14")
	case 0xf0:
		fmt.Println("15")
	default:
		fmt.Println("oopsie poopsie, decode didn't work")
	}
}

func (mem *EmuMemory) Decrement() {
	for mem.delayTimer > 0 {
		time.Sleep(time.Microsecond * 16666)
		mem.delayTimer--
		fmt.Printf("delay timer is now %d \n", mem.delayTimer)
	}
}

func (mem *EmuMemory) SetDelayTimer(time uint8) {
	mem.delayTimer = time
}
