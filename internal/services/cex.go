package services

import (
	"fmt"
	"slices"
	"sync"

	"github.com/mksmstpck/minora-scanner/internal/events/cex"
	"github.com/mksmstpck/minora-scanner/internal/models"
	"github.com/sirupsen/logrus"
)

var blacklist = []string{
	"TRUMPUSDT",
	"ALLUSDT",
	"NEIROUSDT",
}

type Filtered struct {
	Coin   string
	Prices []models.Price
}

func (s *Services) ScanAllExchanges() ([]Filtered, error) {
	exchanges := []struct {
		exchange models.CexType
		api      cex.Cexer
	}{
		{models.Binance, s.events.Cex.Binance},
		{models.Bybit, s.events.Cex.Bybit},
		{models.Gate, s.events.Cex.Gate},
		{models.Kucoin, s.events.Cex.Kucoin},
		{models.Mexc, s.events.Cex.Mexc},
	}

	grouped := make(map[string][]models.Price)

	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, len(exchanges))

	for _, ex := range exchanges {
		wg.Add(1)
		go func(cex models.CexType, api cex.Cexer) {
			defer wg.Done()

			coins, err := api.GetFuturesTicker()
			if err != nil {
				errCh <- fmt.Errorf("error from %d: %w", ex, err)
				return
			}

			for _, coin := range coins {
				mu.Lock()
				grouped[coin.Symbol] = append(grouped[coin.Symbol], models.Price{
					CexType: cex,
					Price:   coin.Price,
				})
				mu.Unlock()
			}
		}(ex.exchange, ex.api)
	}

	wg.Wait()
	close(errCh)

	for e := range errCh {
		fmt.Println("Scan error:", e)
	}

	var filtered []Filtered
	for coin, prices := range grouped {
		if len(prices) > 1 {
			filtered = append(filtered, Filtered{
				Coin:   coin,
				Prices: prices,
			})
		}
	}

	return filtered, nil
}

func (s *Services) SeekPairs() ([]models.Pair, error) {
	filtered, err := s.ScanAllExchanges()
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	var pairs []models.Pair

	for _, f := range filtered {
		if len(f.Prices) == 0 {
			continue
		}

		max := f.Prices[0]
		min := f.Prices[0]

		for _, p := range f.Prices {
			if p.Price > max.Price {
				max = p
			}
			if p.Price < min.Price {
				min = p
			}
		}

		if min.Price == 0 {
			continue
		}

		spread := ((max.Price - min.Price) / min.Price) * 100.0
		if spread > 5 {
			pairs = append(pairs, models.Pair{
				PriceHigh:      max,
				PriceLow:       min,
				Coin:           f.Coin,
				SpreadPercents: spread,
			})
		}
	}

	var newPairs []models.Pair
	for _, pair := range pairs {
		if !s.storage.CheckExists(pair) && !slices.Contains(blacklist, pair.Coin) {
			newPairs = append(newPairs, pair)
			s.storage.SetPair(pair)
		}
	}

	return newPairs, nil
}
