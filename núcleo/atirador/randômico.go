package atirador

import (
	"math/rand"
	"sync"
	"time"
)

var origemRandômica *rand.Rand

func init() {
	origemRandômica = rand.New(&travaOrigemRandômica{
		Source: rand.NewSource(time.Now().UnixNano()),
	})
}

type travaOrigemRandômica struct {
	sync.Mutex
	rand.Source
}

func (t *travaOrigemRandômica) Int63() int64 {
	t.Lock()
	defer t.Unlock()
	return t.Source.Int63()
}

func (t *travaOrigemRandômica) Seed(seed int64) {
	t.Lock()
	defer t.Unlock()
	t.Source.Seed(seed)
}
