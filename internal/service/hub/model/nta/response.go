package nta

type Response struct {
	Data   any    `json:"data"`
	Cursor string `json:"cursor,omitempty"`
}

type TypedResponse[T any] struct {
	Data   T      `json:"data"`
	Cursor string `json:"cursor,omitempty"`
}
