package constant

import "fmt"

func BuildServiceName(service string) string {
	return fmt.Sprintf("global-indexer.%s", service)
}
