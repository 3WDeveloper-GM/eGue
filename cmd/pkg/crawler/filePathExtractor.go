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

type filepathExtractor struct {
	processedAmount int
	log             FileLogger
}

func NewFilepathExtractor(logger FileLogger) *filepathExtractor {
	return &filepathExtractor{log: logger, processedAmount: 0}
}

// filepathExtractor.Process generates a set of file paths to valid email files located in
// the mailDir directory inside the mails directory. These file paths are later used to locate
// each email to be processed later in the pipeline.
func (fe *filepathExtractor) Process(ctx context.Context, p pipes.Payload) (pipes.Payload, error) {
	payload := p.(*crawlerPayload)
	var files strings.Builder

	fe.log.Log(fmt.Sprintf("processing in %s", payload.rootdir))
	err := filepath.WalkDir(payload.rootdir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if !d.IsDir() {
			fe.processedAmount++
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

	successMessage := fmt.Sprintf("extracted from %s, total file names: %d", payload.rootdir, fe.processedAmount)
	fe.log.Log(successMessage)
	return payload, nil
}
