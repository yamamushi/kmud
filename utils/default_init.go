package utils

import (
	"math/rand"
	"runtime"
	"time"
)

func DefaultInit() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UnixNano())
	go InterruptSignalHandler()
}
