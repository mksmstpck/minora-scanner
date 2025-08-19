package models

import "strings"

type KucoinStandartyzer struct{}

func (s *KucoinStandartyzer) Standartize(rawSymbol string) string {
	return strings.TrimSuffix(rawSymbol, "M")
}
