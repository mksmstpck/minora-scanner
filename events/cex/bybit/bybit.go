package bybit

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/mksmstpck/minora-scanner/config"
	"github.com/mksmstpck/minora-scanner/events/cex"
)

type Bybit struct {
	config config.Config
	client *http.Client
}

func NewBybit(config config.Config, client *http.Client) cex.Cexer {
	return &Bybit{
		config: config,
		client: client,
	}
}

type ticker struct {
	result result
}

type result struct {
	List []listItem `json:"list"`
}

type listItem struct {
	Symbol    string `json:"symbol"`
	LastPrice string `json:"last_price"`
}

func (b *Bybit) Get_futures_ticker() ([]cex.CexResultListItem, error) {
	resp, err := b.client.Get(b.config.BybitUrl)
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
	for i := len(result.result.List); i >= 0; i++ {

		price, err := strconv.ParseFloat(result.result.List[i].LastPrice, 64)
		if err != nil {
			log.Printf("an error occured %s", err)
			return nil, err
		}

		cexResult = append(
			cexResult,
			cex.CexResultListItem{
				Symbol: result.result.List[i].Symbol,
				Price:  price,
			})
	}

	return cexResult, nil
}
