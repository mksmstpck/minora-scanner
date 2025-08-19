package models

import (
	"errors"
)

type CexType int

const (
	Binance CexType = iota
	Bybit
	Gate
	Kucoin
	Mexc
)

type Price struct {
	CexType CexType
	Price   float64
}

type Pair struct {
	PriceHigh      Price
	PriceLow       Price
	Coin           string
	SpreadPercents float64
}

type SymbolStandartizer interface {
	Standartize(rawSymbol string) string
}

type RawSymbol struct {
	Symbol  string
	CexType CexType
}

var standartizers = map[CexType]SymbolStandartizer{
	Binance: &BinanceStandartyzer{},
	Bybit:   &BybitStandartyzer{},
	Gate:    &GateStandartyzer{},
	Kucoin:  &KucoinStandartyzer{},
	Mexc:    &MexcStandartyzer{},
}

func GetStandartizedSymbol(rs RawSymbol) (string, error) {
	standartizer, ok := standartizers[rs.CexType]
	if !ok {
		return "", errors.New("no standartizer found")
	}

	return standartizer.Standartize(rs.Symbol), nil
}
