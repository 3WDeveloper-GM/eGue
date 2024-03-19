package indexer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
)

type Response struct {
	Page struct {
		PageNum  int `json:"page_num"`
		PageSize int `json:"page_size"`
		Total    int `json:"total"`
	} `json:"page"`
}

func IndexMails() {
	log.Println("Starting to index")
	config := &Config{
		PayloadPath:   "./data/payloadFile",
		APIstring:     "http://localhost:4080/api/_bulk",
		Admin:         "admin",
		AdminPassword: "Complexpass#123",
		DataPath:      "filepathFolder",
	}

	CPUprofiler, err := os.Create("CpuPerf.prof")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer CPUprofiler.Close()

	pprof.StartCPUProfile(CPUprofiler)
	defer pprof.StopCPUProfile()

	MemProfiler, err := os.Create("MemPerf.prof")
	if err != nil {
		log.Fatal(err)
		return
	}

	pprof.WriteHeapProfile(MemProfiler)
	defer MemProfiler.Close()

	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}

	counter, err := FileCrawl(rootDir, config.DataPath, config)
	if err != nil {
		log.Fatal(err)
		return
	}

	// // seding final payload
	// err = SendPayload(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = FlushPayloadContents(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Printf("Finished crawling %d emails", counter)
}

func StartIndexing() {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:4080/api/index", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.StatusCode)
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	responseMap := &Response{}
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		log.Fatal(errors.New(err.Error() + "hmmm"))
	}

	if responseMap.Page.Total == 0 {
		IndexMails()
	} else {
		log.Println("Nothing to Index")
	}

}
