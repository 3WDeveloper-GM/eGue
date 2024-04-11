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

// Clone is returns a cloned instance of a pipes.Payload object, this
// is done in order to do Parallel operations with multiple processors
// in 1-N Broadcast stages as implemented in the pipes.pipeline package.
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

// MaskAsProcessed is a function that returns a pipes.Payload object
// to the payload pool. This ensures that the objects are collected by
// go's GC in a more efficient manner.
func (p *crawlerPayload) MarkAsProcessed() {
	p.filePath = p.filePath[:0]
	p.RawContent.Reset()
	p.emailBluePrint = nil
	payloadPool.Put(p)
}
