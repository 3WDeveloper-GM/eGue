package zs

// SearchRequest is used to marshal a client request into
// a ZincSearch compatible request after validation.
type SearchRequest struct {
	Type   string `json:"type"`
	Query  string `json:"query"`
	Field  string `json:"field"`
	From   int    `json:"from"`
	MaxRes int    `json:"max_res"`
}
