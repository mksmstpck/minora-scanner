package gate

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/mksmstpck/minora-scanner/config"
	"github.com/mksmstpck/minora-scanner/events/cex"
)

type Gate struct {
	config config.Config
	client *http.Client
}

func NewBybit(config config.Config, client *http.Client) cex.Cexer {
	return &Gate{
		config: config,
		client: client,
	}
}

type ticker struct {
	Contract string `json:"contract"`
	Last     string `json:"last"`
}

func (b *Gate) Get_futures_ticker() ([]cex.CexResultListItem, error) {
	resp, err := b.client.Get(b.config.GateUrl)
	if err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	var result []ticker
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	var cexResult []cex.CexResultListItem
	for i := len(result); i >= 0; i++ {

		price, err := strconv.ParseFloat(result[i].Last, 64)
		if err != nil {
			log.Printf("an error occured %s", err)
			return nil, err
		}

		cexResult = append(
			cexResult,
			cex.CexResultListItem{
				Symbol: result[i].Contract,
				Price:  price,
			})
	}

	return cexResult, nil
}
