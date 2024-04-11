package domain

import (
	"fmt"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter/zs"
	Validator "github.com/3WDeveloper-GM/pipeline/cmd/pkg/domain/validator"
)

func ValidateInput(v *Validator.Validator, input *zs.SearchRequest) bool {

	section := "search_type"
	notAllowed := "%s not in allowed, allowed search types %v"
	emptymessage := "must not be empty"
	searchtypes := []string{"matchphrase", "match", "term", "querystring", "prefix", "fuzzy", "wildcard"}

	v.Check(v.In(input.Type, searchtypes...), section, fmt.Sprintf(notAllowed, input.Type, searchtypes))
	v.Check(input.Type != "", section, emptymessage)

	section = "query"
	largerMessage := "query must be larger than %d bytes"
	lessthanMessage := "query must be less than %d bytes"

	v.Check(len(input.Query) != 0, section, emptymessage)
	v.Check(len(input.Query) > 1, section, fmt.Sprintf(largerMessage, 1))
	v.Check(len(input.Query) <= 200, section, fmt.Sprintf(lessthanMessage, 200))

	section = "field"
	fieldTypes := []string{"To", "From", "Subject", "CC", "Bcc", "_all"}

	v.Check(input.Field != "", section, emptymessage)
	v.Check(v.In(input.Field, fieldTypes...), section, fmt.Sprintf(notAllowed, input.Field, fieldTypes))

	section = "from"

	mustBePositive := "must be a positive integer"
	//v.Check(input.From == 0, section, emptymessage)
	v.Check(input.From >= 0, section, mustBePositive)

	section = "max_results"
	mustBeLessThan := "%s must be less than %d"

	v.Check(input.MaxRes != 0, section, emptymessage)
	v.Check(input.MaxRes > 0, section, mustBePositive)
	v.Check(input.MaxRes < 200, section, fmt.Sprintf(mustBeLessThan, section, 200))

	return v.Valid()
}
