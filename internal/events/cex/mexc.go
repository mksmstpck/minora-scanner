package cex

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mksmstpck/minora-scanner/internal/config"
	"github.com/mksmstpck/minora-scanner/internal/models"
	"github.com/sirupsen/logrus"
)

var _ Cexer = (*Mexc)(nil)

type Mexc struct {
	config config.Config
	client *http.Client
}

func NewMexc(config config.Config, client *http.Client) *Mexc {
	return &Mexc{
		config: config,
		client: client,
	}
}

type tickerMexc struct {
	Data []dataMexc `json:"data"`
}

type dataMexc struct {
	Symbol    string  `json:"symbol"`
	LastPrice float64 `json:"lastPrice"`
}

func (b *Mexc) GetFuturesTicker() ([]CexResultListItem, error) {
	resp, err := b.client.Get(b.config.MexcUrl)
	if err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	var result tickerMexc
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	var cexResult []CexResultListItem
	for _, r := range result.Data {
		if r.LastPrice == 0 {
			logrus.Printf("empty price string for %s", r.Symbol)
			continue
		}

		symbol, err := models.GetStandartizedSymbol(
			models.RawSymbol{Symbol: r.Symbol, CexType: models.Mexc},
		)
		if err != nil {
			return nil, err
		}

		cexResult = append(
			cexResult,
			CexResultListItem{
				Symbol: symbol,
				Price:  r.LastPrice,
			})
	}

	return cexResult, nil
}
