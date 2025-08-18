package kucoin

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/mksmstpck/minora-scanner/config"
	"github.com/mksmstpck/minora-scanner/events/cex"
)

type Kucoin struct {
	config config.Config
	client *http.Client
}

func NewBybit(config config.Config, client *http.Client) cex.Cexer {
	return &Kucoin{
		config: config,
		client: client,
	}
}

type ticker struct {
	Data []data `json:"data"`
}

type data struct {
	Symbol string `json:"symbol"`
	Price  string `json:"pirce"`
}

func (b *Kucoin) GetFuturesTicker() ([]cex.CexResultListItem, error) {
	resp, err := b.client.Get(b.config.KucoinUrl)
	if err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	var result ticker
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	var cexResult []cex.CexResultListItem
	for i := len(result.Data); i >= 0; i++ {

		price, err := strconv.ParseFloat(result.Data[i].Price, 64)
		if err != nil {
			log.Printf("an error occured %s", err)
			return nil, err
		}

		cexResult = append(
			cexResult,
			cex.CexResultListItem{
				Symbol: result.Data[i].Symbol,
				Price:  price,
			})
	}

	return cexResult, nil
}
