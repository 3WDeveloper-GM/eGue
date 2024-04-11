package crawler

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler/pipes"
)

type FileSource struct {
	mailDirs     []string
	currentDir   string
	currentIndex int
	log          FileLogger
}

func NewFileSource(logger FileLogger) *FileSource {
	return &FileSource{log: logger}
}

func (f *FileSource) Next(ctx context.Context) bool {
	if f.currentIndex >= len(f.mailDirs) {
		return false
	}
	f.currentDir = f.mailDirs[f.currentIndex]
	f.currentIndex++
	return true
}

func (f *FileSource) Payload() pipes.Payload {
	crawlerP := payloadPool.Get().(*crawlerPayload)
	crawlerP.rootdir = f.currentDir
	crawlerP.emailBluePrint = map[string]string{
		"Message-ID":                "",
		"Date":                      "",
		"From":                      "",
		"To":                        "",
		"Subject":                   "",
		"Mime-Version":              "",
		"Content-Type":              "",
		"Content-Transfer-Encoding": "",
		"X-From":                    "",
		"X-To":                      "",
		"X-cc":                      "",
		"X-bcc":                     "",
		"X-Folder":                  "",
		"X-Origin":                  "",
		"X-FileName":                "",
		"Cc":                        "",
		"Bcc":                       "",
		"Body":                      "",
	}
	return crawlerP
}

func (f *FileSource) Error() error {
	return nil
}

func (f *FileSource) ListDirs(path string) error {
	var filenames strings.Builder
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// worst case scenario, the slice is the same length as
	// the directories in the path

	for _, file := range files {
		if file.IsDir() {
			filenames.WriteString(path + "/" + file.Name() + " ")
		}
	}

	f.mailDirs = strings.Fields(filenames.String())
	f.log.Log(fmt.Sprintf("file amount %d", len(f.mailDirs)))

	return nil
}

func (f *FileSource) PrintDirs() {
	for _, dir := range f.mailDirs {
		log.Println(dir)
	}
}
