package coingecko

import (
	"net/http"

	"github.com/mksmstpck/minora-scanner/internal/config"
)

type Coingecko struct {
	config config.Config
	client *http.Client
}

func NewCoingecko(config config.Config, client *http.Client) *Coingecko {
	return &Coingecko{
		config: config,
		client: client,
	}
}
