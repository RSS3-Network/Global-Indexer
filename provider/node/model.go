package node

import (
	"github.com/ethereum/go-ethereum/common"
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
	Request        int
	InvalidRequest int
}

type Cache struct {
	Address  string `json:"address"`
	Endpoint string `json:"endpoint"`
}

type ActivitiesResponse struct {
	Data []*Feed     `json:"data"`
	Meta *MetaCursor `json:"meta,omitempty"`
}

type MetaCursor struct {
	Cursor string `json:"cursor"`
}

type Feed struct {
	ID       string `json:"id"`
	Owner    string `json:"owner,omitempty"`
	Network  string `json:"network"`
	From     string `json:"from"`
	To       string `json:"to"`
	Tag      string `json:"tag"`
	Type     string `json:"type"`
	Platform string `json:"platform,omitempty"`
}
