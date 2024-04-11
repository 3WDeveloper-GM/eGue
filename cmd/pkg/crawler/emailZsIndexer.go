package crawler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler/pipes"
)

// mailIndexer holds the fields necessary in order to send payloads
// to ZinkSearch and index them correctly.
type mailIndexer struct {
	cfg        adapter.DBImplementation
	postClient PostClient
	sentItems  int
	logger     FileLogger
}

func NewMailIndexer(cfg adapter.DBImplementation, client PostClient, logger FileLogger) *mailIndexer {

	return &mailIndexer{
		cfg:        cfg,
		postClient: client,
		sentItems:  0,
		logger:     logger,
	}
}

// mailIndexerProcess indexes a block of emails from a certain user in the mailDir directory.
// It generates a request that will be sent to the /_bulk endpoint in the ZIncSearch Database
// for bulk ingestion of data.
func (mi *mailIndexer) Process(ctx context.Context, p pipes.Payload) (pipes.Payload, error) {
	payload := p.(*crawlerPayload)

	mi.logger.Log(fmt.Sprintf("sending %d bytes to ZS", len(payload.RawContent.Bytes())))

	url := mi.cfg.GetDBURL() + "_bulk"
	req, err := mi.postClient.Generate(http.MethodPost, url, &payload.RawContent)
	if err != nil {
		return nil, err
	}

	//req.Header.Set("Content-Encoding", "gzip")

	resp, err := mi.postClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	mi.logger.Log(resp.Status)

	if resp.StatusCode == 500 {
		body, _ := io.ReadAll(resp.Body)
		switch {
		case string(body) == "{\"error\":\"bufio.Scanner: token too long\"}":
			log.Printf("%s long file, payload size %d bytes", payload.rootdir, len(payload.RawContent.Bytes()))

			var buffered bytes.Buffer

			buffered.Write([]byte(payload.rootdir))
			io.Copy(&buffered, &payload.RawContent)
			buffered.Write([]byte("\n\n\n"))

			file, err := os.OpenFile("logged", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			file.Write(buffered.Bytes())

		default:
			log.Println(string(body))
		}
		return nil, nil
	}

	// if the response is not an ok, or any other, just discard the
	// payload, it can be indexed later.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("we´ve got an error: %s", string(body))
		return nil, nil
	}

	mi.sentItems++

	mi.logger.Log(fmt.Sprintf("discarding payload, %d payloads sent", mi.sentItems))
	return payload, nil
}
