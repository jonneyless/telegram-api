package telegram_api

import (
	"sync"
	"testing"
)

func TestConcurrency(t *testing.T) {
	NewTelegramApi("test:token")
	GetTelegramApi()

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			GetTelegram(1)
			GetTelegramApi()
		})
	}
	wg.Wait()
}
