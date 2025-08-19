package cex

import (
	"net/http"

	"github.com/mksmstpck/minora-scanner/internal/config"
	"github.com/mksmstpck/minora-scanner/internal/models"
)

type Cex struct {
	config  config.Config
	Binance *Binance
	Bybit   *Bybit
	Gate    *Gate
	Kucoin  *Kucoin
	Mexc    *Mexc
}

func NewCex(config config.Config, client *http.Client) *Cex {
	return &Cex{
		config:  config,
		Binance: NewBinance(config, client),
		Bybit:   NewBybit(config, client),
		Gate:    NewGate(config, client),
		Kucoin:  NewKucoin(config, client),
		Mexc:    NewMexc(config, client),
	}
}

type Cexer interface {
	GetFuturesTicker() ([]CexResultListItem, error)
}

type CexResultListItem struct {
	Symbol string
	Price  float64
}

type UntouchedSymbol struct {
	CexType models.CexType
	Symbol  string
}

type CexRusultListItemStandart interface {
	Standart(symbol UntouchedSymbol)
}
