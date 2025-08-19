package models

type BinanceStandartyzer struct{}

func (s *BinanceStandartyzer) Standartize(rawSymbol string) string {
	return rawSymbol
}
