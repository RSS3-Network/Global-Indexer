package rules

import (
	"regexp"
)

type RUCalculator func(uri string) int64

const (
	prefixData   = "data"
	prefixSearch = "search"
)

var (
	Prefix2RUCalculator = map[string]RUCalculator{
		prefixData:   dataRUCalculator,
		prefixSearch: searchRUCalculator,
	}
)

func calculateRU(uri string, ruReMap map[*regexp.Regexp]int64) int64 {
	for re, ru := range ruReMap {
		if re.MatchString(uri) {
			return ru
		}
	}

	return 0
}
