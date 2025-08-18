package binance

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/mksmstpck/minora-scanner/config"
	"github.com/mksmstpck/minora-scanner/events/cex"
)

type Binance struct {
	config config.Config
	client *http.Client
}

func NewBinance(config config.Config, client *http.Client) cex.Cexer {
	return &Binance{
		config: config,
		client: client,
	}
}

type ticker struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

func (b *Binance) Get_futures_ticker() ([]cex.CexResultListItem, error) {
	resp, err := b.client.Get(b.config.BinanceUrl)
	if err != nil {
		log.Printf("an error occured %s", err)
	}

	defer resp.Body.Close()

	var result []ticker
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	var cexResult []cex.CexResultListItem
	for i := len(result); i >= 0; i++ {

		price, err := strconv.ParseFloat(result[i].Price, 64)
		if err != nil {
			log.Printf("an error occured %s", err)
			return nil, err
		}

		cexResult = append(
			cexResult,
			cex.CexResultListItem{
				Symbol: result[i].Symbol,
				Price:  price,
			})
	}

	return cexResult, nil
}
