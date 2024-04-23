package dsl

// ActivitiesRequest represents the request for activities by an account.
type ActivitiesRequest struct {
	Account        string   `param:"account"`
	Limit          *int     `query:"limit" default:"100" min:"1" max:"100"`
	ActionLimit    *int     `query:"action_limit" default:"10" min:"1" max:"20"`
	Cursor         *string  `query:"cursor"`
	SinceTimestamp *uint64  `query:"since_timestamp"`
	UntilTimestamp *uint64  `query:"until_timestamp"`
	Status         *bool    `query:"success"`
	Direction      *string  `query:"direction"`
	Network        []string `query:"network"`
	Tag            []string `query:"tag"`
	Type           []string `query:"-"`
	Platform       []string `query:"platform"`
}

// ActivityRequest represents the request for an activity by its ID.
type ActivityRequest struct {
	ID          string `param:"id"`
	ActionLimit int    `query:"action_limit" default:"10" min:"1" max:"20"`
	ActionPage  int    `query:"action_page" default:"1" min:"1"`
}
