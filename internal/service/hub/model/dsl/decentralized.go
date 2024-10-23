package dsl

// ActivityRequest represents the request for an activity by its ID.
type ActivityRequest struct {
	ID          string `param:"id" validate:"required"`
	ActionLimit int    `query:"action_limit" validate:"min=1,max=20" default:"10"`
	ActionPage  int    `query:"action_page" validate:"min=1" default:"1"`
}

// ActivitiesRequest represents the request for activities by an account.
type ActivitiesRequest struct {
	Account        string   `param:"account" validate:"required"`
	Limit          *int     `query:"limit" validate:"min=1,max=100" default:"100"`
	ActionLimit    *int     `query:"action_limit" validate:"min=1,max=20" default:"10"`
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

// AccountsActivitiesRequest represents the request for activities by multiple accounts.
type AccountsActivitiesRequest struct {
	Accounts       []string `json:"accounts" validate:"required,max=20"`
	Limit          int      `json:"limit" validate:"min=1,max=100" default:"100"`
	ActionLimit    int      `json:"action_limit" validate:"min=1,max=20" default:"10"`
	Cursor         *string  `json:"cursor"`
	SinceTimestamp *uint64  `json:"since_timestamp"`
	UntilTimestamp *uint64  `json:"until_timestamp"`
	Status         *bool    `json:"success"`
	Direction      *string  `json:"direction"`
	Network        []string `json:"network"`
	Tag            []string `json:"tag"`
	Type           []string `json:"type"`
	Platform       []string `json:"platform"`
}

// NetworkActivitiesRequest represents the request for activities by a network.
type NetworkActivitiesRequest struct {
	Network string `param:"network" validate:"required"`

	Limit          int      `query:"limit" validate:"min=1,max=100" default:"100"`
	ActionLimit    int      `query:"action_limit" validate:"min=1,max=20" default:"10"`
	Cursor         *string  `query:"cursor"`
	SinceTimestamp *uint64  `query:"since_timestamp"`
	UntilTimestamp *uint64  `query:"until_timestamp"`
	Status         *bool    `query:"success"`
	Direction      *string  `query:"direction"`
	Tag            []string `query:"tag"`
	Type           []string `query:"-"`
	Platform       []string `query:"platform"`
}

// PlatformActivitiesRequest represents the request for activities by a platform.
type PlatformActivitiesRequest struct {
	Platform string `param:"platform" validate:"required"`

	Limit          int      `query:"limit" validate:"min=1,max=100" default:"50"`
	ActionLimit    int      `query:"action_limit" validate:"min=1,max=20" default:"10"`
	Cursor         *string  `query:"cursor"`
	SinceTimestamp *uint64  `query:"since_timestamp"`
	UntilTimestamp *uint64  `query:"until_timestamp"`
	Status         *bool    `query:"success"`
	Direction      *string  `query:"direction"`
	Tag            []string `query:"tag"`
	Type           []string `query:"-"`
	Network        []string `query:"network"`
}
