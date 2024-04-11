package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/adapter/zs"
	"github.com/3WDeveloper-GM/pipeline/cmd/pkg/domain"
)

type SearchMapper struct {
	adapter adapter.SearchAdapter
	cfg     adapter.DBImplementation
	input   zs.SearchRequest
}

func NewSearchMapper(cfg adapter.DBImplementation, adapter adapter.SearchAdapter) *SearchMapper {
	return &SearchMapper{
		adapter: adapter,
		cfg:     cfg,
	}
}

func (sm *SearchMapper) SetInput(input zs.SearchRequest) {
	sm.input = input
}

func (sm *SearchMapper) Search(index string) ([]domain.Email, error) {

	var response zs.SearchResponse
	// Mapping the request into a datbase searchable form
	// using a map is more convenient than using a string
	// with this method is easier to create a request that
	// zincSearch can understand

	//log.Println(sm)

	requestBody := map[string]interface{}{
		"search_type": sm.input.Type,
		"query": map[string]interface{}{
			"term":  sm.input.Query,
			"field": sm.input.Field,
		},
		"from":        sm.input.From,
		"max_results": sm.input.MaxRes,
		"source":      make([]string, 1),
	}

	jsonBytes, err := json.Marshal(&requestBody)
	if err != nil {
		return nil, err
	}

	//log.Println(string(jsonBytes))

	url := sm.cfg.GetDBURL() + index + "/_search"

	//log.Println(url)

	req, err := sm.adapter.Generate(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}

	resp, err := sm.adapter.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	//log.Println(response)

	return sm.mapMails(&response)
}

func (sm *SearchMapper) mapMails(response *zs.SearchResponse) ([]domain.Email, error) {
	var mails = make([]domain.Email, len(response.Hits.Hits))

	for index, item := range response.Hits.Hits {
		var mail domain.Email

		//getting the response struct into a byte form
		mailInBytes, err := json.Marshal(item.Source)
		if err != nil {
			return nil, err
		}

		// Unmarshalling into a mail form, the mail will lose
		// all the X-related fields from the database.
		err = json.Unmarshal(mailInBytes, &mail)
		if err != nil {
			return nil, err
		}

		mails[index] = mail
	}

	return mails, nil
}
