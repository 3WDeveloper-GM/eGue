package crawler

import (
	"bytes"
	"sync"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler/pipes"
)

var (
	_ pipes.Payload = (*crawlerPayload)(nil)

	payloadPool = sync.Pool{
		New: func() any { return new(crawlerPayload) },
	}
)

type crawlerPayload struct {
	filePath []string
	rootdir  string

	emailBluePrint     map[string]string
	processedmailcount int
	RawContent         bytes.Buffer
}

func (p *crawlerPayload) Clone() pipes.Payload {
	// cloning all the fields into a new crawlerPayload
	newPayload := payloadPool.Get().(*crawlerPayload)
	newPayload.filePath = p.filePath
	newPayload.emailBluePrint = map[string]string{
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
	newPayload.RawContent = p.RawContent
	return newPayload
}

func (p *crawlerPayload) MarkAsProcessed() {
	p.filePath = p.filePath[:0]
	p.RawContent.Reset()
	p.emailBluePrint = nil
	payloadPool.Put(p)
}
