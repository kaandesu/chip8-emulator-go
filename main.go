package main

import "sync"

type emulator struct {
	memory          [4096]byte
	registers       [18]byte
	pc              uint16
	stack           [32]uint16
	sp              uint16
	I               uint16
	delayTimer      uint8
	soundTimer      uint8
	delayTimerMutex sync.Mutex
	soundTimerMutex sync.Mutex
}

func main() {
}
