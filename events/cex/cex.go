package cex

import "github.com/mksmstpck/minora-scanner/config"

type Cex struct {
	config config.Config
}

func NewCex(config config.Config) *Cex {
	return &Cex{
		config: config,
	}
}

type CexResultListItem struct {
	Symbol string
	Price  float64
}

type Cexer interface {
	Get_futures_ticker() ([]CexResultListItem, error)
}
