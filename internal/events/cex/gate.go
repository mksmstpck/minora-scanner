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

var _ Cexer = (*Gate)(nil)

type Gate struct {
	config config.Config
	client *http.Client
}

func NewGate(config config.Config, client *http.Client) *Gate {
	return &Gate{
		config: config,
		client: client,
	}
}

type tickerGate struct {
	Contract string `json:"contract"`
	Last     string `json:"last"`
}

func (g *Gate) GetFuturesTicker() ([]CexResultListItem, error) {
	activeSymbols, err := g.getCexInfo()
	if err != nil {
		logrus.Errorf("an error occured %s", err)
		return nil, err
	}

	resp, err := g.client.Get(g.config.GateTickerUrl)
	if err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	var result []tickerGate
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("an error occured %s", err)
		return nil, err
	}

	var cexResult []CexResultListItem
	for _, r := range result {
		if !activeSymbols[r.Contract] {
			continue
		}

		if r.Last == "" {
			logrus.Printf("empty price string for %s", r.Contract)
			continue
		}

		price, err := strconv.ParseFloat(r.Last, 64)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		symbol, err := models.GetStandartizedSymbol(
			models.RawSymbol{Symbol: r.Contract, CexType: models.Gate},
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

type cexInfoListItemGate struct {
	Symbol string `json:"name"`
	Status string `json:"status"`
}

func (g *Gate) getCexInfo() (map[string]bool, error) {
	resp, err := g.client.Get(g.config.GateCexInfoUrl)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	var result []cexInfoListItemGate
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logrus.Error(err)
		return nil, err
	}

	validated := make(map[string]bool)
	for _, symbol := range result {
		if symbol.Status == "trading" {
			validated[symbol.Symbol] = true
		}
	}

	return validated, nil
}
