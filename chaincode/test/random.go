package test

import (
	"math/rand"
	"time"
)

func NewTimeRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
