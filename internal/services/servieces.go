package services

import (
	"github.com/mksmstpck/minora-scanner/internal/config"
	"github.com/mksmstpck/minora-scanner/internal/events"
	"github.com/mksmstpck/minora-scanner/internal/models"
	"github.com/mksmstpck/minora-scanner/internal/storage"
)

type Services struct {
	events  events.Events
	storage storage.Storager
	config  config.Config
}

type Servicer interface {
	SeekPairs() ([]models.Pair, error)
	ScanAllExchanges() ([]Filtered, error)
}

func NewServiecs(
	events events.Events,
	storage storage.Storager,
	config config.Config,
) Servicer {
	return &Services{
		events:  events,
		storage: storage,
		config:  config,
	}
}
