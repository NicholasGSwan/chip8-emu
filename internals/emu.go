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

var (
	Memory         [4096]byte
	pixels         [32][64]bool
	programCounter uint16
	indexRegister  uint16
	stack          []uint16
	delayTimer     timer
	soundTimer     timer
	varRegister    [16]byte
)

type opCode struct {
	opval    byte
	nib2     byte
	nib3     byte
	nib4     byte
	threeVal uint16
	twoVal   byte
}

type timer uint8

const (
	pcStartPoint = 512
)

var SetXtoYInShiftInstruction = true

var keypad = input.GetKeyPad()

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

func init() {
	copy(Memory[:], Font)
	//mem.stack = make([]uint16, 0)
	programCounter = pcStartPoint
	// display.init()

}

func RunEmu(filepath string) {

	//init()
	running := true
	LoadRom(filepath)
	go runDelayTimer()

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {

			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if k, ok := keypad.ScanToKey[e.Keysym.Scancode]; ok {
					if e.State == sdl.PRESSED {
						k.State = input.ON
					} else {
						k.State = input.OFF
					}
				} else if e.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
			}
		}

		opc := Fetch()
		Decode(opc)
		// time.Sleep(time.Microsecond * 16666)
		// bufio.NewReader(os.Stdin).ReadString('\n')
	}
}

func LoadRom(filepath string) {
	file, err := os.Open(filepath)

	if err != nil {
		fmt.Printf("Could not open file: %s\n", err.Error())
	} else {
		numBytes, err := file.Read(Memory[512:])
		if err != nil {
			fmt.Printf("Could not load rom into memory: %s\n", err.Error())
		}

		fmt.Printf("Number of bytes loaded: %d\n", numBytes)
	}
	fmt.Println("Rom loaded!")
}

func FetchAt(memAddr uint16) opCode {
	op1 := Memory[memAddr]
	op2 := Memory[memAddr+1]

	opcode := opCode{op1 & 0xF0 >> 4, op1 & 0x0F, op2 & 0xF0 >> 4, op2 & 0x0F, (uint16(op1)&0x000F)<<8 + uint16(op2), op2}
	//mem.programCounter = mem.programCounter + 2
	return opcode
}

// Fetches the opcode at the PC, then increments it to the next position
func Fetch() opCode {
	opcode := FetchAt(programCounter)
	incrementPc()
	return opcode
}

func incrementPc() {
	programCounter += 2
}

func decrementPc() {
	programCounter -= 2
}

// decode the various opcodes and execute
func Decode(opcode opCode) {
	fmt.Printf("The Program Counter is: %d\n", programCounter)
	// nib2 := op1 & 0x0F
	switch opcode.opval {
	case 0x00:
		switch opcode.twoVal {
		case 0xE0:
			//clear screen
			fmt.Println("clearing screen")
		case 0xEE:
			programCounter = popStack()
			//mem.incrementPc()
		}

		fmt.Println("0")
	case 0x1:
		//jump to threeval
		// fmt.Printf("Program counter is: %d and three val is : %d\n", mem.programCounter, opcode.threeVal)
		// if mem.programCounter-2 == opcode.threeVal {
		// 	os.Exit(0)
		// }
		programCounter = opcode.threeVal

		//fmt.Println("1")
	case 0x2:

		stack = append(stack, programCounter)
		programCounter = opcode.threeVal
		fmt.Println("2")

	case 0x3:
		fmt.Println("3")
		//skip if variable register at nib2 equals twoval
		if varRegister[opcode.nib2] == opcode.twoVal {
			incrementPc()
		}

	case 0x4:
		fmt.Println("4")
		if varRegister[opcode.nib2] != opcode.twoVal {
			incrementPc()
		}
	case 0x5:
		fmt.Println("5")
		if varRegister[opcode.nib2] == varRegister[opcode.nib3] {
			incrementPc()
		}
	case 0x6:
		fmt.Println("6")
		//set variable register at nib2 addr
		varRegister[opcode.nib2] = opcode.twoVal
	case 0x7:
		fmt.Println("7")
		//add value to variable register at nib2
		varRegister[opcode.nib2] += opcode.twoVal
	case 0x8:
		fmt.Println("8")
		switch opcode.nib4 {
		case 0:
			varRegister[opcode.nib2] = varRegister[opcode.nib3]
		case 1:
			varRegister[opcode.nib2] = varRegister[opcode.nib2] | varRegister[opcode.nib3]
		case 2:
			varRegister[opcode.nib2] = varRegister[opcode.nib2] & varRegister[opcode.nib3]
		case 3:
			varRegister[opcode.nib2] = varRegister[opcode.nib2] ^ varRegister[opcode.nib3]
		case 4:
			varRegister[opcode.nib2] = varRegister[opcode.nib2] + varRegister[opcode.nib3]
		case 5:
			varRegister[opcode.nib2] = varRegister[opcode.nib2] - varRegister[opcode.nib3]
		case 6:
			//shift bit right and put in VF
			if SetXtoYInShiftInstruction {
				varRegister[opcode.nib2] = varRegister[opcode.nib3]
			}

			if bit := 0x1 & varRegister[opcode.nib2]; bit == 0x1 {
				varRegister[0xF] = 1
			} else {
				varRegister[0xF] = 0
			}

			varRegister[opcode.nib2] = varRegister[opcode.nib2] >> 1
		case 7:
			varRegister[opcode.nib2] = varRegister[opcode.nib3] - varRegister[opcode.nib2]
		case 0xE:
			//shift bit left and put in VF
			if SetXtoYInShiftInstruction {
				varRegister[opcode.nib2] = varRegister[opcode.nib3]
			}

			if bit := 0x8 & varRegister[opcode.nib2]; bit == 0x8 {
				varRegister[0xF] = 1
			} else {
				varRegister[0xF] = 0
			}

			varRegister[opcode.nib2] = varRegister[opcode.nib2] << 1
		}
	case 0x9:
		fmt.Println("9")
		if varRegister[opcode.nib2] != varRegister[opcode.nib3] {
			incrementPc()
		}
	case 0xA:
		fmt.Println("10")
		//set index register
		//fmt.Printf("setting index register to %d\n", opcode.threeVal)
		indexRegister = opcode.threeVal
	case 0xB:
		fmt.Println("11")
		//jump with offset.  Has ambiguous behavior; implementing original functionality (BNNN vs BXNN)
		programCounter = uint16(varRegister[0]) + opcode.threeVal
	case 0xc:
		fmt.Println("12")
		//Random
		varRegister[opcode.nib2] = byte(rand.Int()) & opcode.twoVal
	case 0xd:
		fmt.Println("13")
		// bitwise AND behaves the same as modulo for powers of 2 apparently
		x := varRegister[opcode.nib2] & 63
		y := varRegister[opcode.nib3] & 31
		fmt.Printf("x is : %d and y is : %d\n", x, y)
		fmt.Printf("start is : %d and end is : %d\n", indexRegister, indexRegister+uint16(opcode.nib4))
		spriteData := Memory[indexRegister : indexRegister+uint16(opcode.nib4)]
		fmt.Printf("The sprite data is %d in length\n", len(spriteData))
		drawDisplay(x, y, spriteData)
		time.Sleep(time.Microsecond * 66664)
		// time.Sleep(time.Millisecond * 332)
		//display/draw
	case 0xe:
		fmt.Println("14")
		//skip instruction if key is/ is not pressed

		key := keypad.ScanToKey[keypad.ValToScan[opcode.nib2]]

		switch opcode.twoVal {
		case 0x9e:
			if key.State {
				incrementPc()
			}
		case 0xa1:
			if !key.State {
				incrementPc()
			}

		}
	case 0xf:
		fmt.Println("15")
		//wait for keypress

		switch opcode.twoVal {
		case 0x07:
			//sets VX to the current value of the delay timer
			varRegister[opcode.nib2] = uint8(delayTimer)
		case 0x15:
			//sets the delay timer to the value in VX
			delayTimer.setTimer(varRegister[opcode.nib2])
		case 0x18:
			soundTimer.setTimer(varRegister[opcode.nib2])
		case 0x33:
			num := varRegister[opcode.nib2]
			storeDigits(num)
			//mem.getDigits()
		case 0x55:
			for i := 0; i <= int(opcode.nib2); i++ {
				Memory[indexRegister+uint16(i)] = varRegister[i]
			}
		case 0x65:
			for i := 0; i <= int(opcode.nib2); i++ {
				varRegister[i] = Memory[indexRegister+uint16(i)]
			}

		case 0x1e:
			//The index register I will get the value in VX added to it.
			tmp := indexRegister
			indexRegister += uint16(varRegister[opcode.nib2])
			if indexRegister < tmp {
				varRegister[0xf] = 1
			}
		case 0x0a:
			keyIsPressed := false
			var keyval byte
			for _, v := range keypad.ScanToKey {
				if v.State {
					keyIsPressed = bool(v.State)
					keyval = v.Value
				}
			}
			if !keyIsPressed {
				decrementPc()
			} else {
				varRegister[opcode.nib2] = keyval
			}
		}

	default:
		fmt.Println("oopsie poopsie, decode didn't work")
	}

}

func drawDisplay(x, y byte, spriteData []byte) {
	varRegister[0xf] = 0

	for i := 0; i < len(spriteData) && int(y)+i < 32; i++ {
		bRow := spriteData[i]
		xb := x
		bs := 7
		for xb < 64 && xb < x+8 {
			if (bRow & (1 << bs)) != 0 {
				if pixels[int(y)+i][xb] {
					pixels[int(y)+i][xb] = false
					varRegister[0xf] = 1
				} else {
					pixels[int(y)+i][xb] = true
				}
			}
			bs--
			xb++
		}
	}
	display.DrawDisplay(pixels)
	// mem.printDisplay()
}

// original way I was simulating the display rendering before implementing sdl2
func printDisplay() {
	display.DrawDisplay(pixels)
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[0]); x++ {
			if pixels[y][x] {

				fmt.Print("11")

			} else {
				fmt.Print("00")
			}
		}
		fmt.Print("\n")
	}

	fmt.Print("\n\n\n")

}

func (t timer) Decrement() {
	if t > 0 {
		t--
		time.Sleep(time.Microsecond * 16666)
		fmt.Printf("delay timer is now %d \n", t)
	}
}

func SetDelayTimer(time uint8) {
	delayTimer = timer(time)
}
func (t timer) setTimer(time uint8) {
	t = timer(time)
}

func runDelayTimer() {
	for {
		for delayTimer > 0 {
			delayTimer.Decrement()
		}
	}
}

func runSoundTimer() {
	for {
		for soundTimer > 0 {
			soundTimer.Decrement()
		}
	}
}

func popStack() uint16 {

	retVal := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return retVal
}

func (opcode opCode) GetThreeval() uint16 {
	return opcode.threeVal
}

func PrintOpCode(opcode opCode) {
	fmt.Printf("Opcode vals are:\n opval: %x\nnib2: %x\nnib3: %x\nnib4: %x\nthreeVal: %x\ntwoVal: %x\n", opcode.opval, opcode.nib2, opcode.nib3, opcode.nib4, opcode.threeVal, opcode.twoVal)
}

func PrintPcCounter() {
	fmt.Printf("Program Counter is at: %d\n", programCounter)
}

func storeDigits(num byte) {
	ir := indexRegister + 2
	fmt.Printf("preparing to store %d into memory\n", num)

	for i := 0; i < 3; i++ {
		fmt.Printf("Putting %d into memory addr %d\n", num%10, ir)
		Memory[ir] = num % 10
		ir--
		num = num / 10
	}
}
func getDigits() {
	fmt.Print("The digits just stored are ")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d", Memory[indexRegister+uint16(i)])
	}
	fmt.Print("\n")
}

func DestroyResources() {
	display.DestroyResources()
}
