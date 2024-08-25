package main

import (
	"log"
	"os"
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

const memoryOffset = 0x200

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
	screen          [64][32]int
}

type settings struct {
	title      string
	width      int32
	height     int32
	fps        int32
	pixelScale int32
}

var (
	Emulator      *emulator
	Settings      settings
	screenImage   *rl.Image
	screenTexture rl.Texture2D
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
	Settings = settings{
		width:      64,
		height:     32,
		fps:        60,
		pixelScale: 8,
		title:      "CHIP-8 emulator Go",
	}
	rl.InitWindow(Settings.width*Settings.pixelScale, Settings.height*Settings.pixelScale, Settings.title)
	rl.SetTargetFPS(Settings.fps)
	rl.SetTraceLogLevel(rl.LogError)
	rl.SetExitKey(0)

	Emulator = &emulator{}

	Emulator.sp = 0
	Emulator.pc = 0x200
	err := Emulator.LoadROM("./demos/IBM Logo.ch8")
	if err != nil {
		log.Fatalf("Failed to load ROM: %v", err)
	}

	screenImage = rl.GenImageColor(int(Settings.width), int(Settings.height), rl.DarkGray)
	screenTexture = rl.LoadTextureFromImage(screenImage)
}

func input() {}

func render() {
	rl.BeginDrawing()
	Emulator.execute()
	rl.UpdateTexture(screenTexture, rl.LoadImageColors(screenImage))
	rl.DrawTextureEx(screenTexture, rl.NewVector2(0, 0), 0, float32(Settings.pixelScale), rl.White)
	rl.EndDrawing()
}

func quit() {
	rl.CloseWindow()
	rl.UnloadTexture(screenTexture)
	rl.UnloadImage(screenImage)
}

func (e *emulator) LoadROM(filename string) error {
	rom, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	for i := 0; i < len(rom) && i+memoryOffset < len(e.memory); i++ {
		e.memory[i+memoryOffset] = rom[i]
	}

	return nil
}
