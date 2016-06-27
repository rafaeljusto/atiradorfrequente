package randômico

import (
	"math/rand"
	"sync"
	"time"
)

// FonteRandômica fonte de geração de números aleatórios livre de problemas de
// concorrência.
var FonteRandômica *rand.Rand

func init() {
	FonteRandômica = rand.New(&travaFonteRandômica{
		Source: rand.NewSource(time.Now().UnixNano()),
	})
}

// travaFonteRandômica segue a estratégia definida em
// https://nishanths.svbtle.com/do-not-seed-the-global-random
type travaFonteRandômica struct {
	sync.Mutex
	rand.Source
}

func (t *travaFonteRandômica) Int63() int64 {
	t.Lock()
	defer t.Unlock()
	return t.Source.Int63()
}

func (t *travaFonteRandômica) Seed(seed int64) {
	t.Lock()
	defer t.Unlock()
	t.Source.Seed(seed)
}
