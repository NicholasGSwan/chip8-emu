package emu

import (
	"fmt"
	"time"
)

type EmuMemory struct {
	Memory         [4096]byte
	programCounter uint16
	indexRegister  uint16
	stack          []uint16
	delayTimer     uint8
	soundTimer     uint8
	varRegister    [16]byte
}

type opCode struct {
	opval    byte
	nib2     byte
	nib3     byte
	nib4     byte
	threeVal uint16
	twoVal   byte
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
	copy(mem.Memory[:], font)
	//mem.stack = make([]uint16, 0)
	mem.programCounter = 500
}

func (mem *EmuMemory) FetchAt(memAddr uint16) opCode {
	op1 := mem.Memory[memAddr]
	op2 := mem.Memory[memAddr+1]

	opcode := opCode{op1 & 0xF0 >> 4, op1 & 0x0F, op2 & 0xF0 >> 4, op2 & 0x0F, (uint16(op1)&0x000F)<<8 + uint16(op2), op2}
	//mem.programCounter = mem.programCounter + 2
	return opcode
}

func (mem *EmuMemory) Fetch() opCode {
	opcode := mem.FetchAt(mem.programCounter)
	mem.incrementPc()
	return opcode
}

func (mem *EmuMemory) incrementPc() {
	mem.programCounter += 2
}

func (mem *EmuMemory) Decode(opcode opCode) {

	// nib2 := op1 & 0x0F
	switch opcode.opval {
	case 0x00:
		switch opcode.twoVal {
		case 0xE0:
			//clear screen
			fmt.Println("clearing screen")
		case 0xEE:
			mem.programCounter = mem.popStack()
			//mem.incrementPc()
		}

		fmt.Println("0")
	case 0x10:
		//jump to threeval
		mem.programCounter = opcode.threeVal

		fmt.Println("1")
	case 0x20:

		mem.stack = append(mem.stack, mem.programCounter)
		mem.programCounter = opcode.threeVal
		fmt.Println("2")

	case 0x30:
		fmt.Println("3")
		//skip if variable register at nib2 equals twoval
		if mem.varRegister[opcode.nib2] == opcode.twoVal {
			mem.incrementPc()
		}

	case 0x40:
		fmt.Println("4")
		if mem.varRegister[opcode.nib2] != opcode.twoVal {
			mem.incrementPc()
		}
	case 0x50:
		fmt.Println("5")
		if mem.varRegister[opcode.nib2] == mem.varRegister[opcode.nib3] {
			mem.incrementPc()
		}
	case 0x60:
		fmt.Println("6")
		//set variable register at nib2 addr
		mem.varRegister[opcode.nib2] = opcode.twoVal
	case 0x70:
		fmt.Println("7")
		//add value to variable register at nib2
		mem.varRegister[opcode.nib2] += opcode.twoVal
	case 0x80:
		fmt.Println("8")
		switch opcode.nib4 {
		case 0:
			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib3]
		case 1:
			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib2] | mem.varRegister[opcode.nib3]
		case 2:
			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib2] & mem.varRegister[opcode.nib3]
		case 3:
			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib2] ^ mem.varRegister[opcode.nib3]
		case 4:
			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib2] + mem.varRegister[opcode.nib3]
		case 5:
			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib2] - mem.varRegister[opcode.nib3]
		case 7:
			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib3] - mem.varRegister[opcode.nib2]
		}
	case 0x90:
		fmt.Println("9")
		if mem.varRegister[opcode.nib2] != mem.varRegister[opcode.nib3] {
			mem.incrementPc()
		}
	case 0xA0:
		fmt.Println("10")
		//set index register
		mem.indexRegister = opcode.threeVal
	case 0xB0:
		fmt.Println("11")
	case 0xc0:
		fmt.Println("12")
	case 0xd0:
		fmt.Println("13")
		//display/draw
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

func (mem *EmuMemory) popStack() uint16 {

	retVal := mem.stack[len(mem.stack)-1]
	mem.stack = mem.stack[:len(mem.stack)-1]
	return retVal
}

func PrintOpCode(opcode opCode) {
	fmt.Printf("Opcode vals are:\n opval: %x\nnib2: %x\nnib3: %x\nnib4: %x\nthreeVal: %x\ntwoVal: %x\n", opcode.opval, opcode.nib2, opcode.nib3, opcode.nib4, opcode.threeVal, opcode.twoVal)
}
