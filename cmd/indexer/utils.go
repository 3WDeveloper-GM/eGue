package indexer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

type Config struct {
	DataPath      string
	PayloadPath   string
	APIstring     string
	Admin         string
	AdminPassword string
}

// var wg sync.WaitGroup
var eg errgroup.Group

func FileCrawl(root, path string, config *Config) (int64, error) {

	var counter int64 = 0
	var chunkSize int64 = 10000
	var blockSize int = 20

	err := filepath.WalkDir(fmt.Sprintf("%s/%s", root, path), func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if (counter+1)%chunkSize == 0 {

			if err := eg.Wait(); err != nil {
				log.Fatal(err)
				return err
			}
			log.Println("Sending...")
			for index := 0; index < blockSize; index++ {
				eg.Go(func() error {
					//defer wg.Done()
					err := SendPayload(config, index)
					if err != nil {
						log.Fatal(err)
						return err
					}
					//log.Println("Flushing...")
					err = FlushPayloadContents(config, index)
					if err != nil {
						log.Fatal(err)
						return err
					}
					return nil
				})

			}
			if err := eg.Wait(); err != nil {
				log.Fatal(err)
				return err
			}
		}

		if !dirEntry.IsDir() {
			completePath := fmt.Sprintf(config.PayloadPath+"%d"+".ndjson", counter%int64(blockSize))
			eg.Go(func() error {
				return ConstructPayload(path, completePath)
			})
			counter++
		}
		return err
	})
	if err != nil {
		return counter, err
	}

	//Final check for sending the data and flushing the remainder files in the data directory
	err = filepath.WalkDir(fmt.Sprintf("%s/%s", root, "data"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
			return err
		}

		if !d.IsDir() {
			err = SendFinalPayload(config, path)
			if err != nil {
				log.Fatal(err)
				return err
			}
			err = FlushPayloadFolderContents(path)
			if err != nil {
				log.Fatal(err)
				return err
			}
		}
		return nil
	})

	return counter, err
}

func ConstructPayload(filePath, payloadPath string) error {

	//defer wg.Done()

	keyValuesArray := map[string]string{
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
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var emailbody strings.Builder
	var currentline string
	var lineErr error

	separatorFound := false
	currentHeader := ""
	for lineErr == nil {
		currentline, lineErr = reader.ReadString('\n')
		//fmt.Println(currentline)

		if currentline == "\r\n" {
			separatorFound = true
		}

		if !separatorFound {

			headervalues := strings.Split(currentline, ":")
			if len(headervalues) == 1 {
				headervalues = append(headervalues, " ")
			}

			_, ok := keyValuesArray[headervalues[0]]
			if ok {
				keyValuesArray[headervalues[0]] = headervalues[1]
				currentHeader = headervalues[0]
			} else {
				keyValuesArray[currentHeader] += " " + currentline
			}

		} else {
			_, err := emailbody.WriteString(currentline)
			if err != nil {
				return err
			}
		}

	}

	textItem := &Email{
		MessageID:              keyValuesArray["Message-ID"],
		Date:                   keyValuesArray["Date"],
		From:                   keyValuesArray["From"],
		To:                     keyValuesArray["To"],
		Subject:                keyValuesArray["Subject"],
		MimeVersion:            keyValuesArray["Mime-Version"],
		ContentType:            keyValuesArray["Content-Type"],
		ContentTranferEncoding: keyValuesArray["Content-Transfer-Encoding"],
		XFrom:                  keyValuesArray["X-From"],
		Xto:                    keyValuesArray["X-To"],
		Xcc:                    keyValuesArray["X-cc"],
		Xbcc:                   keyValuesArray["X-bcc"],
		XFolder:                keyValuesArray["X-Folder"],
		XOrigin:                keyValuesArray["X-Origin"],
		XFileName:              keyValuesArray["X-FileName"],
		CC:                     keyValuesArray["Cc"],
		BCC:                    keyValuesArray["Bcc"],
		Bod:                    emailbody.String(),
	}

	jsonArr, err := json.Marshal(textItem)
	if err != nil {
		return err
	}

	prependString := "{ \"index\" : {\"_index\": \"mails\" }}\n"

	jsonArr = append([]byte(prependString), jsonArr...)
	jsonArr = append(jsonArr, []byte{'\n'}...)

	file2, err := os.OpenFile(payloadPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file2.Close()

	_, err = file2.Write(jsonArr)
	if err != nil {
		return err
	}

	return nil
}

func SendPayload(config *Config, index int) error {
	completePath := fmt.Sprintf(config.PayloadPath+"%d"+".ndjson", index)
	payload, err := os.OpenFile(completePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer payload.Close()

	req, err := http.NewRequest(http.MethodPost, config.APIstring, payload)
	if err != nil {
		log.Fatal(err)
		return err
	}

	req.SetBasicAuth(config.Admin, config.AdminPassword)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()
	//log.Println(resp.Status)

	return nil
}

func FlushPayloadContents(config *Config, index int) error {
	completePath := fmt.Sprintf(config.PayloadPath+"%d"+".ndjson", index)
	err := os.Remove(completePath)
	if err != nil {
		return err
	}
	return nil
}

func SendFinalPayload(config *Config, path string) error {
	payload, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)

	if err != nil {
		log.Fatal(err)
		return err
	}
	defer payload.Close()

	req, err := http.NewRequest(http.MethodPost, config.APIstring, payload)
	if err != nil {
		log.Fatal(err)
		return err
	}

	req.SetBasicAuth(config.Admin, config.AdminPassword)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	return nil
}

func FlushPayloadFolderContents(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}
