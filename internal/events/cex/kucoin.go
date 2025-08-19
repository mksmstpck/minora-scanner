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

func (b *Kucoin) GetFuturesTicker() ([]CexResultListItem, error) {
	resp, err := b.client.Get(b.config.KucoinUrl)
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
