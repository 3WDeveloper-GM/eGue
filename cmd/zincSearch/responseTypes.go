package zincsearch

type Search struct {
	Type       string   `json:"type"`
	SearchTerm string   `json:"search_term"`
	Field      string   `json:"field"`
	SortFields []string `json:"sort_fields"`
	From       int      `json:"from"`
	MaxResults int      `json:"max_results"`
	Source     []string `json:"_source,omitempty"`
}
