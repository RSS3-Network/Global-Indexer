package rules

import (
	"regexp"
)

var searchRuReMap = map[*regexp.Regexp]int64{
	regexp.MustCompile(`^/v2/recent-activities.*$`): 10,
	regexp.MustCompile(`^/suggestions/.*$`):         2,
	regexp.MustCompile(`^/dapps.*$`):                2,
	regexp.MustCompile(`^/activities.*$`):           10,
	regexp.MustCompile(`^/v2/activities.*$`):        10,
	regexp.MustCompile(`^/activities/.*$`):          1,
	regexp.MustCompile(`^/v2/activities/.*$`):       1,
}

func searchRUCalculator(uri string) int64 {
	return calculateRU(uri, searchRuReMap)
}
