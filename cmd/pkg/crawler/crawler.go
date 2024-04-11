package crawler

import (
	"context"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler/pipes"
)

// Crawler is a public struct that exposes a pipes.Pipeline object
// configured in a certain way, the configuration details are found in
// the assemblePipeline function.
type Crawler struct {
	p *pipes.Pipeline
}

func NewCrawler(cfg adapter.DBImplementation, client *adapter.SearchAdapter, logger FileLogger) *Crawler {
	return &Crawler{p: assemblePipeline(cfg, client, logger)}
}

// assemblePipeline is a private function that creates a pipes.pipeline object with a
// set configuration, this configuration consists of three worker pools with different sizes.
// This function only needs to be changed whenever the developer wants to try other processing
// schemes, like FIFO, or 1-N Broadcast as implemented in the pipes package.
func assemblePipeline(cfg adapter.DBImplementation, client *adapter.SearchAdapter, logger FileLogger) *pipes.Pipeline {
	return pipes.New(
		pipes.DynamicWorkerPool(NewFilepathExtractor(logger), 16),
		pipes.DynamicWorkerPool(NewMailProcessor(logger), 64),
		pipes.DynamicWorkerPool(NewMailIndexer(cfg, client, logger), 128),
	)
}

// Crawl starts the processing directives enabled by the pipes.Pipeline object with
// the provided source.
func (c *Crawler) Crawl(ctx context.Context, source FileSource) (int, error) {
	sink := new(countingSink)
	err := c.p.Process(ctx, &source, sink)
	return sink.getCount(), err
}
