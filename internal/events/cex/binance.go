package cex

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mksmstpck/minora-scanner/internal/config"
	"github.com/mksmstpck/minora-scanner/internal/models"
	"github.com/sirupsen/logrus"
)

var _ Cexer = (*Binance)(nil)

type Binance struct {
	config config.Config
	client *http.Client
}

func NewBinance(config config.Config, client *http.Client) *Binance {
	return &Binance{
		config: config,
		client: client,
	}
}

type tickerBinance struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (b *Binance) GetFuturesTicker() ([]CexResultListItem, error) {
	activeSymbols, err := b.getCexInfo()
	if err != nil {
		logrus.Errorf("failed to get CEX info: %s", err)
		return nil, err
	}

	resp, err := b.client.Get(b.config.BinanceTickerUrl)
	if err != nil {
		logrus.Errorf("an error occured %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result []tickerBinance
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logrus.Errorf("an error occured %s", err)
		return nil, err
	}

	var cexResult []CexResultListItem
	for _, r := range result {
		if !activeSymbols[r.Symbol] {
			continue
		}

		if r.Price == "" {
			logrus.Infof("empty price string for %s", r.Symbol)
			continue
		}
		price, err := strconv.ParseFloat(r.Price, 64)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		symbol, err := models.GetStandartizedSymbol(
			models.RawSymbol{Symbol: r.Symbol, CexType: models.Binance},
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

type infoBinance struct {
	Symbols []symbolBinance `json:"symbols"`
}

type symbolBinance struct {
	Symbol string `json:"symbol"`
	Status string `json:"status"`
}

func (b *Binance) getCexInfo() (map[string]bool, error) {
	resp, err := b.client.Get(b.config.BinanceCexInfoUrl)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	var result infoBinance
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logrus.Error(err)
		return nil, err
	}

	validated := make(map[string]bool) // Initialize the map
	for _, symbol := range result.Symbols {
		if symbol.Status == "TRADING" {
			validated[symbol.Symbol] = true
		}
	}

	return validated, nil
}
