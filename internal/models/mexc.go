package models

import "strings"

type MexcStandartyzer struct{}

func (s *MexcStandartyzer) Standartize(rawSymbol string) string {
	return strings.ReplaceAll(rawSymbol, "_", "")
}
