package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/3WDeveloper-GM/pipeline/cmd/app"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/logger"
	indexer "github.com/3WDeveloper-GM/pipeline/cmd/pkg/service/Indexer"
)

// This file is for compiling the indexer as a standalone tool in the app.
// With this we can estimate how the tool performs on its own.

func main() {

	config := app.NewZsConfiguration()
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 64
	transport.MaxConnsPerHost = 128
	transport.MaxIdleConnsPerHost = 128

	client := &http.Client{
		Transport: transport,
	}

	logfile, err := os.Create("logfile.log")
	if err != nil {
		panic(err)
	}

	logger := logger.NewLogger(logfile)
	adapter := adapter.NewAdapter(client, config)

	index := indexer.NewMailIndexer(config, adapter, logger)

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()

	root, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	err = index.Index(ctx, root)
	if err != nil {
		log.Println(err)
		return
	}
}
