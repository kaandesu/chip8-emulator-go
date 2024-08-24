package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

/*
NOTE:
Read the instruction that PC is currently pointing at from memory.
An instruction is two bytes, so you will need to read two successive
bytes from memory and combine them into one 16-bit instruction.
You should then immediately increment the PC by 2, to be ready to
fetch the next opcode. Some people do this during the “execute”
stage, since some instructions will increment it by 2 more to skip an
instruction, but in my opinion that’s very error-prone.
Code duplication is a bad thing. If you forget to increment it in
one of the instructions, you’ll have problems. Do it here!
*/
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
	_ = X
	_ = NN
	_ = NNN

	switch inst {
	case 0x0:
		switch Y {
		case 0xE:
			switch N {
			case 0x0:
				// TODO: clear screen
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
		// TODO: draw screen here
		e.screen[e.registers[X]%64][e.registers[Y]%32] = 1
		e.draw(screenImage)

	}
}

func (e *emulator) draw(screenImage *rl.Image) {
	for y := range e.screen {
		for _, x := range e.screen[y] {
			if x > 0 {
				rl.ImageDrawPixel(screenImage, int32(x), int32(int(Settings.height)-y-1), rl.White)
			}
		}
	}
}
