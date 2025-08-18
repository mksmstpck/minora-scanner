package events

import (
	"github.com/mksmstpck/minora-scanner/events/cex"
	"github.com/mksmstpck/minora-scanner/events/coingecko"
)

type Events struct {
	Coingecko *coingecko.Coingecko
	Cex       *cex.Cex
}

func NewEvents() Events {
	return Events{
		Coingecko: coingecko.NewCoingecko(),
		Cex:       cex.NewCex(),
	}
}
