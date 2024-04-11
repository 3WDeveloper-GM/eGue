package zs

//	less is more, with only these parameters
//
// it can search confortably in the ZS engine
type SearchRequest struct {
	Type   string `json:"type"`
	Query  string `json:"query"`
	Field  string `json:"field"`
	From   int    `json:"from"`
	MaxRes int    `json:"max_res"`
}
