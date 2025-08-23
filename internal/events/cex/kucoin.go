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

var _ Cexer = (*Kucoin)(nil)

type Kucoin struct {
	config config.Config
	client *http.Client
}

func NewKucoin(config config.Config, client *http.Client) *Kucoin {
	return &Kucoin{
		config: config,
		client: client,
	}
}

type tickerKucoin struct {
	Data []dataKucoin `json:"data"`
}

type dataKucoin struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (k *Kucoin) GetFuturesTicker() ([]CexResultListItem, error) {
	activeSymbols, err := k.getCexInfo()
	if err != nil {
		logrus.Errorf("an error occured %s", err)
		return nil, err
	}

	resp, err := k.client.Get(k.config.KucoinTickerUrl)
	if err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	var result tickerKucoin
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	var cexResult []CexResultListItem
	for _, r := range result.Data {
		if !activeSymbols[r.Symbol] {
			continue
		}

		if r.Price == "" {
			logrus.Printf("empty price string for %s", r.Symbol)
			continue
		}

		price, err := strconv.ParseFloat(r.Price, 64)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		symbol, err := models.GetStandartizedSymbol(
			models.RawSymbol{Symbol: r.Symbol, CexType: models.Kucoin},
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

type resultKucoin struct {
	Data []dataListItemKucoin `json:"data"`
}

type dataListItemKucoin struct {
	Symbol string `json:"symbol"`
	Status string `json:"status"`
}

func (k *Kucoin) getCexInfo() (map[string]bool, error) {
	resp, err := k.client.Get(k.config.KucoinCexInfoUrl)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	var result resultKucoin
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logrus.Error(err)
		return nil, err
	}

	validated := make(map[string]bool)
	for _, symbol := range result.Data {
		if symbol.Status == "Open" {
			validated[symbol.Symbol] = true
		}
	}

	return validated, nil
}
