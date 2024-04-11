package crawler

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler/pipes"
)

var minCap = 16

type filepathExtractor struct {
	minCapacity int
	log         FileLogger
}

func NewFilepathExtractor(logger FileLogger) *filepathExtractor {
	return &filepathExtractor{minCapacity: minCap, log: logger}
}

// the file extractor is a processor that extracts the mail files from a certain directory
// given by the payload rootdir field.
func (fe *filepathExtractor) Process(ctx context.Context, p pipes.Payload) (pipes.Payload, error) {
	payload := p.(*crawlerPayload)
	counter := 0
	var files strings.Builder

	fe.log.Log(fmt.Sprintf("processing in %s", payload.rootdir))
	err := filepath.WalkDir(payload.rootdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if !d.IsDir() {
			counter++
			files.WriteString(path + " ")
		}
		return nil
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	payload.filePath = strings.Fields(files.String())

	// if the directory is empty, just remove the payload from the pipeline
	if len(payload.filePath) == 0 {
		fe.log.Log("empty directory")
		return nil, nil
	}

	fe.log.Log(fmt.Sprintf("finished extracting file names from %s, amount of extracted file names: %d", payload.rootdir, counter))
	return payload, nil
}
