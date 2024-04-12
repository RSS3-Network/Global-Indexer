package nta

type Response struct {
	Data   any    `json:"data"`
	Cursor string `json:"cursor,omitempty"`
}
