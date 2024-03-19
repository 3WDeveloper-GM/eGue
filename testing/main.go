package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Response struct {
	List []string `json:"list"`
	Page struct {
		PageNum  int `json:"page_num"`
		PageSize int `json:"page_size"`
		Total    int `json:"total"`
	} `json:"page"`
}

func main() {

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

	fmt.Println(len(responseMap.List))

	file, err := os.Create("./response.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer file.Close()

}
