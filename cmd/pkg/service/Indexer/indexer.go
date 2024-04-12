package indexer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/domain"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/logger"
)

var mailsPath = domain.MailDir

// MailIndexerService: is a service for the Index handler in the application
// this struct implements the methods required to make it a DBIndex compatible
// object.
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

// Index: this function implements the Index method from the DBIndex interface
// it checks first if there are entries in the database and if there are no entries
// it doesn't do anything.
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

	err = json.NewDecoder(resp.Body).Decode(&total)
	if err != nil {
		return err
	}

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
