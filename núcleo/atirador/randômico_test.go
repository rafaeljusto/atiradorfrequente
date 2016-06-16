package atirador

import (
	"testing"
	"time"
)

func TestTravaOrigemRandômica_Int63(t *testing.T) {
	iterações := 50
	númeroRandômicoCh := make(chan int, iterações)

	for i := 0; i < iterações; i++ {
		go func() {
			origemRandômica.Seed(time.Now().UnixNano())
			númeroRandômicoCh <- origemRandômica.Intn(100)
		}()
	}

	for i := 0; i < iterações; i++ {
		númeroRandômico := <-númeroRandômicoCh
		if númeroRandômico < 0 || 100 < númeroRandômico {
			t.Errorf("Número randômico %d inexperado", númeroRandômico)
		}
	}
}
