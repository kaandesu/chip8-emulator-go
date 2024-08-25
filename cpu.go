package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (e *emulator) fetch() (b0, b1 uint8) {
	b0 = e.memory[e.pc]
	b1 = e.memory[e.pc+1]
	e.pc += 2
	return
}

func (e *emulator) decode() (inst, X, Y, N, NN uint8, NNN uint16) {
	b0, b1 := e.fetch()
	inst = (b0 & 0xF0) >> 4
	X = b0 & 0x0F
	Y = (b1 & 0xF0) >> 4
	N = b1 & 0x0F
	NN = b1
	NNN = uint16(X)<<8 | uint16(NN)
	return
}

func (e *emulator) execute() {
	inst, X, Y, N, NN, NNN := e.decode()

	switch inst {
	case 0x0:
		switch Y {
		case 0xE:
			switch N {
			case 0x0:
				e.screen = [64][32]int{}
			case 0xE:
				// Ommited
			}
		}
	case 0x1:
		e.pc = NNN
	case 0x6:
		e.registers[X] = NN
	case 0x7:
		e.registers[X] += NN
	case 0xA:
		e.I = NNN
	case 0xD:
		e.drawSprite(e.registers[X], e.registers[Y], N)
	}
}

func (e *emulator) draw(screenImage *rl.Image) {
	for x, row := range e.screen {
		for y, value := range row {
			if value > 0 {
				rl.ImageDrawPixel(screenImage, int32(x), int32(y), rl.White)
			}
		}
	}
}

func (e *emulator) drawSprite(VX, VY, N uint8) {
	x := VX % uint8(Settings.width)
	y := VY % uint8(Settings.height)
	e.registers[0xF] = 0

	for row := uint8(0); row < N; row++ {
		spriteByte := e.memory[e.I+uint16(row)]
		for col := uint8(0); col < 8; col++ {
			if spriteByte&(0x80>>col) != 0 { // Check if the pixel in the sprite is set
				pixelX := (x + col) % uint8(Settings.width)
				pixelY := (y + row) % uint8(Settings.height)

				if e.screen[pixelX][pixelY] == 1 {
					e.registers[0xF] = 1
				}

				e.screen[pixelX][pixelY] ^= 1
			}
		}
	}

	e.draw(screenImage)
}
