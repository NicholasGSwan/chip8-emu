package emu

import (
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/NicholasGSwan/chip8-emu/internals/display"
	"github.com/NicholasGSwan/chip8-emu/internals/input"
	"github.com/veandco/go-sdl2/sdl"
)

type EmuMemory struct {
	Memory         [4096]byte
	display        [32][64]bool
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

const (
	pcStartPoint = 512
)

var keybState []uint8
var SetXtoYInShiftInstruction = true
var kp = input.GetKeys()

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

func (mem *EmuMemory) init() {
	copy(mem.Memory[:], Font)
	//mem.stack = make([]uint16, 0)
	mem.programCounter = pcStartPoint
	// display.init()

}

func (e *EmuMemory) RunEmu(filepath string) {

	// e.Init()
	running := true
	e.LoadRom(filepath)

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {

			case *sdl.QuitEvent:
				running = false
			}
		}

		opc := e.Fetch()
		e.Decode(opc)
		// bufio.NewReader(os.Stdin).ReadString('\n')
	}
}

func (mem *EmuMemory) LoadRom(filepath string) {
	file, err := os.Open(filepath)

	if err != nil {
		fmt.Printf("Could not open file: %s\n", err.Error())
	} else {
		numBytes, err := file.Read(mem.Memory[512:])
		if err != nil {
			fmt.Printf("Could not load rom into memory: %s\n", err.Error())
		}

		fmt.Printf("Number of bytes loaded: %d\n", numBytes)
	}
	fmt.Println("Rom loaded!")
}

func (mem *EmuMemory) FetchAt(memAddr uint16) opCode {
	op1 := mem.Memory[memAddr]
	op2 := mem.Memory[memAddr+1]

	opcode := opCode{op1 & 0xF0 >> 4, op1 & 0x0F, op2 & 0xF0 >> 4, op2 & 0x0F, (uint16(op1)&0x000F)<<8 + uint16(op2), op2}
	//mem.programCounter = mem.programCounter + 2
	return opcode
}

// Fetches the opcode at the PC, then increments it to the next position
func (mem *EmuMemory) Fetch() opCode {
	opcode := mem.FetchAt(mem.programCounter)
	mem.incrementPc()
	return opcode
}

func (mem *EmuMemory) incrementPc() {
	mem.programCounter += 2
}

func (mem *EmuMemory) decrementPc() {
	mem.programCounter -= 2
}

// decode the various opcodes and execute
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
	case 0x1:
		//jump to threeval
		// fmt.Printf("Program counter is: %d and three val is : %d\n", mem.programCounter, opcode.threeVal)
		// if mem.programCounter-2 == opcode.threeVal {
		// 	os.Exit(0)
		// }
		mem.programCounter = opcode.threeVal

		//fmt.Println("1")
	case 0x2:

		mem.stack = append(mem.stack, mem.programCounter)
		mem.programCounter = opcode.threeVal
		fmt.Println("2")

	case 0x3:
		fmt.Println("3")
		//skip if variable register at nib2 equals twoval
		if mem.varRegister[opcode.nib2] == opcode.twoVal {
			mem.incrementPc()
		}

	case 0x4:
		fmt.Println("4")
		if mem.varRegister[opcode.nib2] != opcode.twoVal {
			mem.incrementPc()
		}
	case 0x5:
		fmt.Println("5")
		if mem.varRegister[opcode.nib2] == mem.varRegister[opcode.nib3] {
			mem.incrementPc()
		}
	case 0x6:
		fmt.Println("6")
		//set variable register at nib2 addr
		mem.varRegister[opcode.nib2] = opcode.twoVal
	case 0x7:
		fmt.Println("7")
		//add value to variable register at nib2
		mem.varRegister[opcode.nib2] += opcode.twoVal
	case 0x8:
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
		case 6:
			//shift bit right and put in VF
			if SetXtoYInShiftInstruction {
				mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib3]
			}

			if bit := 0x1 & mem.varRegister[opcode.nib2]; bit == 0x1 {
				mem.varRegister[0xF] = 1
			} else {
				mem.varRegister[0xF] = 0
			}

			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib2] >> 1
		case 7:
			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib3] - mem.varRegister[opcode.nib2]
		case 0xE:
			//shift bit left and put in VF
			if SetXtoYInShiftInstruction {
				mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib3]
			}

			if bit := 0x8 & mem.varRegister[opcode.nib2]; bit == 0x8 {
				mem.varRegister[0xF] = 1
			} else {
				mem.varRegister[0xF] = 0
			}

			mem.varRegister[opcode.nib2] = mem.varRegister[opcode.nib2] << 1
		}
	case 0x9:
		fmt.Println("9")
		if mem.varRegister[opcode.nib2] != mem.varRegister[opcode.nib3] {
			mem.incrementPc()
		}
	case 0xA:
		fmt.Println("10")
		//set index register
		//fmt.Printf("setting index register to %d\n", opcode.threeVal)
		mem.indexRegister = opcode.threeVal
	case 0xB:
		fmt.Println("11")
		//jump with offset.  Has ambiguous behavior; implementing original functionality (BNNN vs BXNN)
		mem.programCounter = uint16(mem.varRegister[0]) + opcode.threeVal
	case 0xc:
		fmt.Println("12")
		//Random
		mem.varRegister[opcode.nib2] = byte(rand.Int()) & opcode.twoVal
	case 0xd:
		fmt.Println("13")
		// bitwise AND behaves the same as modulo for powers of 2 apparently
		x := mem.varRegister[opcode.nib2] & 63
		y := mem.varRegister[opcode.nib3] & 31
		fmt.Printf("x is : %d and y is : %d\n", x, y)
		fmt.Printf("start is : %d and end is : %d\n", mem.indexRegister, mem.indexRegister+uint16(opcode.nib4))
		spriteData := mem.Memory[mem.indexRegister : mem.indexRegister+uint16(opcode.nib4)]
		fmt.Printf("The sprite data is %d in length\n", len(spriteData))
		mem.drawDisplay(x, y, spriteData)

		//display/draw
	case 0xe:
		fmt.Println("14")
		//skip instruction if key is/ is not pressed
		keybState = sdl.GetKeyboardState()
		isPressed := keybState[kp[mem.varRegister[opcode.nib2]]]

		switch opcode.twoVal {
		case 0x9e:
			if isPressed == 1 {
				mem.incrementPc()
			}
		case 0xa1:
			if isPressed == 0 {
				mem.incrementPc()
			}

		}
	case 0xf:
		fmt.Println("15")
		//wait for keypress

		switch opcode.twoVal {
		case 0x07:
			//sets VX to the current value of the delay timer
			mem.varRegister[opcode.nib2] = mem.delayTimer
		case 0x15:
			//sets the delay timer to the value in VX
			mem.delayTimer = mem.varRegister[opcode.nib2]
		case 0x18:
			mem.soundTimer = mem.varRegister[opcode.nib2]
		case 0x33:
			num := mem.varRegister[opcode.nib2]
			mem.storeDigits(num)
			//mem.getDigits()
		case 0x55:
			for i := 0; i <= int(opcode.nib2); i++ {
				mem.Memory[mem.indexRegister+uint16(i)] = mem.varRegister[i]
			}
		case 0x65:
			for i := 0; i <= int(opcode.nib2); i++ {
				mem.varRegister[i] = mem.Memory[mem.indexRegister+uint16(i)]
			}

		case 0x1e:
			//The index register I will get the value in VX added to it.
			tmp := mem.indexRegister
			mem.indexRegister += uint16(mem.varRegister[opcode.nib2])
			if mem.indexRegister < tmp {
				mem.varRegister[0xf] = 1
			}
		case 0x0a:

		}

	default:
		fmt.Println("oopsie poopsie, decode didn't work")
	}
}

func (mem *EmuMemory) drawDisplay(x, y byte, spriteData []byte) {
	mem.varRegister[0xf] = 0

	for i := 0; i < len(spriteData) && int(y)+i < 32; i++ {
		bRow := spriteData[i]
		xb := x
		bs := 7
		for xb < 64 && xb < x+8 {
			if (bRow & (1 << bs)) != 0 {
				if mem.display[int(y)+i][xb] {
					mem.display[int(y)+i][xb] = false
					mem.varRegister[0xf] = 1
				} else {
					mem.display[int(y)+i][xb] = true
				}
			}
			bs--
			xb++
		}
	}
	display.DrawDisplay(mem.display)
	// mem.printDisplay()
}

// original way I was simulating the display rendering before implementing sdl2
func (mem *EmuMemory) printDisplay() {
	display.DrawDisplay(mem.display)
	for y := 0; y < len(mem.display); y++ {
		for x := 0; x < len(mem.display[0]); x++ {
			if mem.display[y][x] {

				fmt.Print("11")

			} else {
				fmt.Print("00")
			}
		}
		fmt.Print("\n")
	}

	fmt.Print("\n\n\n")

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

func (opcode opCode) GetThreeval() uint16 {
	return opcode.threeVal
}

func PrintOpCode(opcode opCode) {
	fmt.Printf("Opcode vals are:\n opval: %x\nnib2: %x\nnib3: %x\nnib4: %x\nthreeVal: %x\ntwoVal: %x\n", opcode.opval, opcode.nib2, opcode.nib3, opcode.nib4, opcode.threeVal, opcode.twoVal)
}

func (mem EmuMemory) PrintPcCounter() {
	fmt.Printf("Program Counter is at: %d\n", mem.programCounter)
}

func (mem *EmuMemory) storeDigits(num byte) {
	ir := mem.indexRegister + 2
	fmt.Printf("preparing to store %d into memory\n", num)

	for i := 0; i < 3; i++ {
		fmt.Printf("Putting %d into memory addr %d\n", num%10, ir)
		mem.Memory[ir] = num % 10
		ir--
		num = num / 10
	}
}
func (mem *EmuMemory) getDigits() {
	fmt.Print("The digits just stored are ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d", mem.Memory[mem.indexRegister+uint16(i)])
	}
	fmt.Print("\n")
}

func DestroyResources() {
	display.DestroyResources()
}
