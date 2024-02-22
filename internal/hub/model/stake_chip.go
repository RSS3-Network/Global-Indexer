package model

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"

	"github.com/naturalselectionlabs/rss3-global-indexer/contract/l2"
)

type StakeChip struct {
	ID       *big.Int        `json:"id"`
	Metadata json.RawMessage `json:"metadata"`
}

func BuildStakeChipMetadata(id *big.Int, metadata json.RawMessage, baseURL url.URL) (json.RawMessage, error) {
	var tokenMetadata l2.ChipsTokenMetadata
	if err := json.Unmarshal(metadata, &tokenMetadata); err != nil {
		return nil, err
	}

	tokenMetadata.Image = baseURL.JoinPath(fmt.Sprintf("/stake/chips/%d/image.svg", id)).String()

	return json.Marshal(tokenMetadata)
}
