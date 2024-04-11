package crawler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/crawler/pipes"
	"golang.org/x/xerrors"
)

// CONSTANTS
// formatting strings for the bytes buffer, these get the correct formatting of
//
//	an .ndjson file for the payload generation, the value of maxSize is the amount
//
// of bytes that the buffer from ZincSearch is able to read without throwing a 500
// error to the client.
var prependString = "{ \"index\" : {\"_index\": \"mails\" }}\n"
var postpendString []byte = []byte{'\n'}
var maxSize = 1024

// var escapePattern = regexp.MustCompile(`(\\u003[0-z])+`)
//var tildePattern = regexp.MustCompile(`\s+`)

// email processor is the pipeline stage that does the email parsing and
// marshalling into a /api/_bulk compatible payload for sending into the
// zincsearch database.
type emailProcessor struct {
	processedCorrectly int
	processedcount     int
	logger             FileLogger
}

func NewMailProcessor(logger FileLogger) *emailProcessor {
	return &emailProcessor{processedcount: 0, logger: logger}
}

func (ep *emailProcessor) Process(ctx context.Context, p pipes.Payload) (pipes.Payload, error) {
	payload := p.(*crawlerPayload)

	//log.Println(len(filesArray))

	var byteContent bytes.Buffer //a byte slice for partitioning the data

	for index, file := range payload.filePath {
		if file == "" {
			continue
		}
		if index > 0 {
			payload.emailBluePrint = map[string]string{
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
		}

		email, err := os.Open(file)
		if err != nil {
			wrapperError := xerrors.Errorf("found some anomaly in the path %s with the index %d : %w", file, index, err)
			return nil, wrapperError
		}
		defer email.Close()

		scanner := bufio.NewScanner(email)
		var emailbody strings.Builder
		var currentLine string

		separatorFound := false
		currentHeader := ""
		for scanner.Scan() {
			currentLine = scanner.Text()

			//removing duplicate whitespace

			currentLine = strings.Join(strings.Fields(currentLine), " ")
			currentLine = strings.ToValidUTF8(currentLine, " ")

			if currentHeader == "X-FileName" {
				separatorFound = true
			}

			if !separatorFound {

				headerValues := strings.Split(currentLine, ":")
				if len(headerValues) == 1 {
					headerValues = append(headerValues, " ")
				}

				headerPrefix := headerValues[0]
				headerSuffix := headerValues[1]

				_, ok := payload.emailBluePrint[headerPrefix]
				if ok {
					payload.emailBluePrint[headerPrefix] = headerSuffix
					currentHeader = headerPrefix
				} else {
					payload.emailBluePrint[currentHeader] += " " + currentLine
				}

			} else {
				_, err := emailbody.WriteString(currentLine)
				if err != nil {
					return nil, err
				}
			}
		}

		payload.emailBluePrint["Body"] = emailbody.String()

		jsonBytes, err := json.Marshal(payload.emailBluePrint)
		if err != nil {
			return nil, err
		}

		// if the size in bytes of the processed mail is larger than 64 kilobytes, it's better to altogether
		// remove the email from the email payload.
		if len(jsonBytes) > (maxSize << (10 * 1)) {
			ep.logger.Log(fmt.Sprintf("mail in %s larger than %dkb", file, maxSize))
			continue
		}

		byteContent.Write([]byte(prependString))
		byteContent.Write(jsonBytes)
		byteContent.Write(postpendString)

		payload.processedmailcount++
		ep.processedcount++
	}

	io.Copy(&payload.RawContent, &byteContent)
	ep.processedCorrectly++
	//log.Println(string(payload.RawContent))
	//log.Println("Finished payload")
	//log.Printf("processed the emails in the directory %s, partition %d", payload.rootdir, partition)

	//messaging to stdin
	var message = "%s processed, %d total files, %d processed correctly"

	ep.logger.Log(fmt.Sprintf(message, payload.rootdir, payload.processedmailcount, ep.processedCorrectly))
	return payload, nil
}
