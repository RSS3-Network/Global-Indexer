package nameresolver

//go:generate go run --mod=mod github.com/dmarkham/enumer --values --type=NameService --linecomment --output name_service_string.go --json --sql
type NameService int

const (
	NameServiceUnknown   NameService = iota // unknown
	NameServiceENS                          // eth
	NameServiceCSB                          // csb
	NameServiceLens                         // lens
	NameServiceFarcaster                    // fc
)
