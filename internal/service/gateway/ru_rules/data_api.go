package rules

import (
	"regexp"
)

var dataRuReMap = map[*regexp.Regexp]int64{
	regexp.MustCompile(`^/accounts/activities.*$`):     10,
	regexp.MustCompile(`^/accounts/.*/activities.*$`):  5,
	regexp.MustCompile(`^/accounts/.*/profiles.*$`):    2,
	regexp.MustCompile(`^/activities/.*$`):             1,
	regexp.MustCompile(`^/mastodon/.*/activities.*$`):  2,
	regexp.MustCompile(`^/networks/.*/activities.*$`):  2,
	regexp.MustCompile(`^/platforms/.*/activities.*$`): 2,
}

func dataRUCalculator(uri string) int64 {
	return calculateRU(uri, dataRuReMap)
}
