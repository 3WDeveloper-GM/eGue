package indexer

import (
	"context"
	"net/http"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/domain"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/logger"
	"github.com/go-chi/render"
)

var mailsPath = domain.MailDir

// this is a blueprint for what an indexer should do
type MailIndexerService struct {
	crawler *crawler.Crawler
	source  *crawler.FileSource
	adapter *adapter.SearchAdapter
	config  adapter.DBImplementation
}

func NewMailIndexer(config adapter.DBImplementation, client *adapter.SearchAdapter, logger *logger.FileLogger) *MailIndexerService {
	return &MailIndexerService{
		crawler: crawler.NewCrawler(config, client, logger),
		source:  crawler.NewFileSource(logger),
		adapter: client,
		config:  config,
	}
}

func (mi *MailIndexerService) Index(ctx context.Context, root string) error {

	url := mi.config.GetDBURL() + "index"
	req, err := mi.adapter.Generate(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := mi.adapter.Do(req)
	if err != nil {
		return err
	}

	var total struct {
		Page struct {
			Total int `json:"total"`
		} `json:"page"`
	}

	render.DecodeJSON(resp.Body, total)

	if total.Page.Total == 0 {
		maildirs := root + mailsPath

		err := mi.source.ListDirs(maildirs)
		if err != nil {
			return err
		}

		_, err = mi.crawler.Crawl(ctx, *mi.source)
		if err != nil {
			return err
		}
	}

	return nil
}
