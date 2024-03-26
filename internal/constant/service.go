package constant

import "fmt"

var ServiceName string

func BuildServiceName() string {
	return fmt.Sprintf("global-indexer.%s", ServiceName)
}
