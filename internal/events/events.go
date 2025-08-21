package events

import (
	"net/http"

	"github.com/mksmstpck/minora-scanner/internal/config"
	"github.com/mksmstpck/minora-scanner/internal/events/cex"
	"github.com/mksmstpck/minora-scanner/internal/events/coingecko"
)

type Events struct {
	Coingecko *coingecko.Coingecko
	Cex       *cex.Cex
}

func NewEvents(config config.Config, client *http.Client) Events {
	return Events{
		Coingecko: coingecko.NewCoingecko(config, client),
		Cex:       cex.NewCex(config, client),
	}
}
