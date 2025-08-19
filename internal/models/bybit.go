package models

type BybitStandartyzer struct{}

func (s *BybitStandartyzer) Standartize(rawSymbol string) string {
	return rawSymbol
}
