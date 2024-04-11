package crawler

import (
	"context"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler/pipes"
)

// this exposes the api for the pipeline package, crawler is
// an implementation of a little file crawler for a technical
// test. It crawls a list of files, processes them and then
// indexes them into the ZIncSearch database
type Crawler struct {
	p *pipes.Pipeline
}

// Is a wrapper that encapsulates all the files and

func NewCrawler(cfg adapter.DBImplementation, client *adapter.SearchAdapter, logger FileLogger) *Crawler {
	return &Crawler{p: assemblePipeline(cfg, client, logger)}
}

func assemblePipeline(cfg adapter.DBImplementation, client *adapter.SearchAdapter, logger FileLogger) *pipes.Pipeline {
	return pipes.New(
		pipes.DynamicWorkerPool(NewFilepathExtractor(logger), 16),
		pipes.DynamicWorkerPool(NewMailProcessor(logger), 64),
		pipes.DynamicWorkerPool(NewMailIndexer(cfg, client, logger), 128),
	)
}

func (c *Crawler) Crawl(ctx context.Context, source FileSource) (int, error) {
	sink := new(countingSink)
	err := c.p.Process(ctx, &source, sink)
	return sink.getCount(), err
}
