package utils

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func Keepliveserver(url string) {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for range ticker.C {
			resp, err := http.Get(url)
			if err != nil {
				log.Error().Msgf("KeepAlive error: %v", err)
				continue
			}
			resp.Body.Close()
		}
	}()
}
