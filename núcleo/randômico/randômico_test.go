package randômico

import (
	"testing"
	"time"
)

func TestTravaFonteRandômica_Int63(t *testing.T) {
	iterações := 50
	númeroRandômicoCh := make(chan int, iterações)

	for i := 0; i < iterações; i++ {
		go func() {
			FonteRandômica.Seed(time.Now().UnixNano())
			númeroRandômicoCh <- FonteRandômica.Intn(100)
		}()
	}

	for i := 0; i < iterações; i++ {
		númeroRandômico := <-númeroRandômicoCh
		if númeroRandômico < 0 || 100 < númeroRandômico {
			t.Errorf("Número randômico %d inexperado", númeroRandômico)
		}
	}
}
