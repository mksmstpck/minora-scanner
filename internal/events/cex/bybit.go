package cex

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/mksmstpck/minora-scanner/internal/config"
	"github.com/mksmstpck/minora-scanner/internal/models"
	"github.com/sirupsen/logrus"
)

var _ Cexer = (*Bybit)(nil)

type Bybit struct {
	config config.Config
	client *http.Client
}

func NewBybit(config config.Config, client *http.Client) *Bybit {
	return &Bybit{
		config: config,
		client: client,
	}
}

type tickerBybit struct {
	result result
}

type result struct {
	List []listItem `json:"list"`
}

type listItem struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"last_price"`
}

func (b *Bybit) GetFuturesTicker() ([]CexResultListItem, error) {
	resp, err := b.client.Get(b.config.BybitUrl)
	if err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	var result tickerBybit
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	var cexResult []CexResultListItem
	for _, r := range result.result.List {
		if r.LastPrice == "" {
			logrus.Printf("empty price string for %s", r.Symbol)
			continue
		}

		price, err := strconv.ParseFloat(r.LastPrice, 64)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		symbol, err := models.GetStandartizedSymbol(
			models.RawSymbol{Symbol: r.Symbol, CexType: models.Bybit},
		)
		if err != nil {
			return nil, err
		}

		cexResult = append(
			cexResult,
			CexResultListItem{
				Symbol: symbol,
				Price:  price,
			})
	}

	return cexResult, nil
}
