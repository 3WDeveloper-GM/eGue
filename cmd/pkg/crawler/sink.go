package crawler

import (
	"context"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler/pipes"
)

type countingSink struct {
	count int
}

func NewCountingSink() *countingSink {
	return &countingSink{count: 1}
}

func (c *countingSink) getCount() int {
	return c.count
}

func (c *countingSink) Consume(_ context.Context, p pipes.Payload) error {
	c.count++
	return nil
}
