package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mksmstpck/minora-scanner/internal/config"
	"github.com/mksmstpck/minora-scanner/internal/services"
	"github.com/sirupsen/logrus"
)

type Handlers struct {
	client *http.Client
	config config.Config
	s      services.Servicer
}

func NewHandlers(client *http.Client, config config.Config, s services.Servicer) *Handlers {
	return &Handlers{
		client: client,
		config: config,
		s:      s,
	}
}

type message struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (h *Handlers) SendPairs() error {
	pairs, err := h.s.SeekPairs()
	if err != nil {
		logrus.Error(err)
		return err
	}

	for _, pair := range pairs {
		text := fmt.Sprintf(
			"Pair:\n Coin - %s\n CexHigh - %s\n CexLow - %s\n Spread - %f",
			pair.Coin,
			pair.PriceHigh.CexType.Name(),
			pair.PriceLow.CexType.Name(),
			pair.SpreadPercents,
		)

		message, err := json.Marshal(message{ChatID: h.config.ChatID, Text: text})
		if err != nil {
			logrus.Error(err)
			return err
		}

		h.client.Post(
			fmt.Sprintf("https://api.telegram.org/bot%s/%s", h.config.Token, "sendMessage"),
			"application/json",
			bytes.NewReader(message),
		)
	}

	return nil
}
