package node

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/protocol-go/schema/filter"
	"github.com/rss3-network/protocol-go/schema/metadata"
)

var (
	RssNodeCacheKey  = "nodes:rss"
	FullNodeCacheKey = "nodes:full"
)

type ActivityRequest struct {
	ID          string `param:"id" description:"Retrieve details for the specified activity ID" examples:"[\"0x5ffa607a127d63fb36827075493d1de06f58fc44710b9ffb887b2effe02d2b8b\"]"`
	ActionLimit int    `query:"action_limit" form:"action_limit" description:"Specify the number of actions within the activity to retrieve" examples:"[10]" default:"10" min:"1" max:"20"`
	ActionPage  int    `query:"action_page" form:"action_page" description:"Specify the pagination for actions" default:"1" min:"1"`
}

type AccountActivitiesRequest struct {
	Account        string   `param:"account" description:"Retrieve activities from the specified account" examples:"[\"vitalik.eth\",\"stani.lens\",\"diygod.csb\"]"`
	Limit          *int     `query:"limit" form:"limit" description:"Specify the number of activities to retrieve" examples:"[20]" default:"100" min:"1" max:"100"`
	ActionLimit    *int     `query:"action_limit" form:"action_limit" description:"Specify the number of actions within the activity to retrieve" examples:"[10]" default:"10" min:"1" max:"20"`
	Cursor         *string  `query:"cursor" form:"cursor" description:"Specify the cursor used for pagination"`
	SinceTimestamp *uint64  `query:"since_timestamp" form:"since_timestamp" description:"Retrieve activities starting from this timestamp" examples:"[1654000000]"`
	UntilTimestamp *uint64  `query:"until_timestamp" form:"until_timestamp" description:"Retrieve activities up to this timestamp" examples:"[1696000000]"`
	Status         *bool    `query:"success" form:"success" description:"Retrieve activities based on status"`
	Direction      *string  `query:"direction" form:"direction" description:"Retrieve activities based on direction"`
	Network        []string `query:"network" form:"network" description:"Retrieve activities from the specified network(s)" examples:"[[\"ethereum\",\"polygon\"]]"`
	Tag            []string `query:"tag" form:"tag" description:"Retrieve activities from the specified tag(s)"`
	Type           []string `query:"" form:"type" description:"Retrieve activities from the specified type(s)"`
	Platform       []string `query:"platform" form:"platform" description:"Retrieve activities from the specified platform(s)"`
}

type DataResponse struct {
	Address        common.Address
	Data           []byte
	First          bool
	Err            error
	Request        int
	InvalidRequest int
}

type ErrResponse struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type Cache struct {
	Address  string `json:"address"`
	Endpoint string `json:"endpoint"`
}

type ActivityResponse struct {
	Data *Feed `json:"data"`
}

type ActivitiesResponse struct {
	Data []*Feed     `json:"data"`
	Meta *MetaCursor `json:"meta,omitempty"`
}

type MetaCursor struct {
	Cursor string `json:"cursor"`
}

type Feed struct {
	ID       string    `json:"id"`
	Owner    string    `json:"owner,omitempty"`
	Network  string    `json:"network"`
	Index    uint      `json:"index"`
	From     string    `json:"from"`
	To       string    `json:"to"`
	Tag      string    `json:"tag"`
	Type     string    `json:"type"`
	Platform string    `json:"platform,omitempty"`
	Actions  []*Action `json:"actions"`
}

type Action struct {
	Tag         string            `json:"tag"`
	Type        string            `json:"type"`
	Platform    string            `json:"platform,omitempty"`
	From        string            `json:"from"`
	To          string            `json:"to"`
	Metadata    metadata.Metadata `json:"metadata"`
	RelatedURLs []string          `json:"related_urls,omitempty"`
}

type Actions []*Action

var _ json.Unmarshaler = (*Action)(nil)

func (a *Action) UnmarshalJSON(bytes []byte) error {
	type ActionAlias Action

	type action struct {
		ActionAlias

		MetadataX json.RawMessage `json:"metadata"`
	}

	var temp action

	err := json.Unmarshal(bytes, &temp)
	if err != nil {
		return fmt.Errorf("unmarshal action: %w", err)
	}

	tag, err := filter.TagString(temp.Tag)
	if err != nil {
		return fmt.Errorf("invalid action tag: %w", err)
	}

	typeX, err := filter.TypeString(tag, temp.Type)
	if err != nil {
		return fmt.Errorf("invalid action type: %w", err)
	}

	temp.Metadata, err = metadata.Unmarshal(typeX, temp.MetadataX)
	if err != nil {
		return fmt.Errorf("invalid action metadata: %w", err)
	}

	*a = Action(temp.ActionAlias)

	return nil
}

// WorkerToNetworksMap Supplement the conditions for a full node based on the configuration file.
// https://github.com/NaturalSelectionLabs/RSS3-Node/blob/develop/deploy/config.development.yaml
var WorkerToNetworksMap = map[filter.Name][]string{
	filter.Fallback:   {filter.NetworkEthereum.String()},
	filter.Mirror:     {filter.NetworkArweave.String()},
	filter.Farcaster:  {filter.NetworkFarcaster.String()},
	filter.RSS3:       {filter.NetworkEthereum.String()},
	filter.Paragraph:  {filter.NetworkArweave.String()},
	filter.OpenSea:    {filter.NetworkEthereum.String()},
	filter.Uniswap:    {filter.NetworkEthereum.String()},
	filter.Optimism:   {filter.NetworkEthereum.String()},
	filter.Aavegotchi: {filter.NetworkPolygon.String()},
	filter.Lens:       {filter.NetworkPolygon.String()},
}

var NetworkToWorkersMap = map[filter.Network][]string{
	filter.NetworkEthereum:  {filter.Fallback.String(), filter.RSS3.String(), filter.OpenSea.String(), filter.Uniswap.String(), filter.Optimism.String()},
	filter.NetworkArweave:   {filter.Mirror.String(), filter.Paragraph.String()},
	filter.NetworkFarcaster: {filter.Farcaster.String()},
	filter.NetworkPolygon:   {filter.Aavegotchi.String(), filter.Lens.String()},
}

var PlatformToWorkerMap = map[filter.Platform]string{
	filter.PlatformRSS3:       filter.RSS3.String(),
	filter.PlatformMirror:     filter.Mirror.String(),
	filter.PlatformFarcaster:  filter.Farcaster.String(),
	filter.PlatformParagraph:  filter.Paragraph.String(),
	filter.PlatformOpenSea:    filter.OpenSea.String(),
	filter.PlatformUniswap:    filter.Uniswap.String(),
	filter.PlatformOptimism:   filter.Optimism.String(),
	filter.PlatformAavegotchi: filter.Aavegotchi.String(),
	filter.PlatformLens:       filter.Lens.String(),
}

var TagToWorkersMap = map[filter.Tag][]string{
	filter.TagTransaction: {filter.Optimism.String()},
	filter.TagCollectible: {filter.OpenSea.String()},
	filter.TagExchange:    {filter.RSS3.String(), filter.Uniswap.String()},
	filter.TagSocial:      {filter.Farcaster.String(), filter.Mirror.String(), filter.Lens.String(), filter.Paragraph.String()},
	filter.TagMetaverse:   {filter.Aavegotchi.String()},
}
