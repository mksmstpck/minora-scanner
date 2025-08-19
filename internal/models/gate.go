package models

import "strings"

type GateStandartyzer struct{}

func (s *GateStandartyzer) Standartize(rawSymbol string) string {
	return strings.ReplaceAll(rawSymbol, "_", "")
}
