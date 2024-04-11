package handlers

import (
	"context"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter/zs"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/domain"
)

type DBSearch interface {
	Search(index string) ([]domain.Email, error)
	SetInput(input zs.SearchRequest)
}

type DBIndex interface {
	Index(ctx context.Context, root string) error
}
