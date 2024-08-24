package main

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

/*NOTE:
  The first CHIP-8 interpreter (on the COSMAC VIP computer)
  was also located in RAM, from address 000 to 1FF. It
  would expect a CHIP-8 program to be loaded into memory after
  it, starting at address 200 (512 in decimal). Although modern
  interpreters are not in the same memory space, you should do the
  same to be able to run the old programs; you can just
  leave the initial space empty, except for the font.
*/

// NOTE: Display: 64 x 32 pixels

type emulator struct {
	memory          [4096]uint8 // 4kb memory
	registers       [16]uint8   // general purpose registers 8-bit
	pc              uint16      // program counter, points to the current instruction
	stack           [32]uint16  // used to call subroutines/functions and return from them
	sp              uint16      // stack pointer starts from -1, set 0 in first use
	I               uint16      // index register which is used to point at locations in memory
	delayTimer      uint8       // decremented at a rate of 60 Hz (60 times per second) until it reaches 0
	soundTimer      uint8       // functions like the delay timer,  also gives off a beeping sound while it’s not 0
	delayTimerMutex sync.Mutex
	soundTimerMutex sync.Mutex
}

type settings struct {
	title  string
	width  int32
	height int32
	fps    int32
}

var (
	Emulator *emulator
	Settings settings
)

func main() {
	setup()
	for !rl.WindowShouldClose() {
		input()
		render()
	}
	defer quit()
}

func setup() {
	Emulator = &emulator{}
	Settings = settings{
		width:  32,
		height: 64,
		fps:    60,
		title:  "CHIP-8 emulator Go",
	}
	// TODO: Load the rom to the 'memory'
	Emulator.sp = 0

	rl.SetTargetFPS(Settings.fps)
	rl.SetTraceLogLevel(rl.LogError)
	rl.InitWindow(Settings.width, Settings.height, Settings.title)
	rl.SetExitKey(0)
}

func input() {}

func render() {}

func quit() {
	rl.CloseWindow()
}
