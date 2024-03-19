package zincsearch

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/3WDeveloper-GM/eGue/cmd/config"
)

func MatchPhrase(app *config.Application, queryOptions *Search) ([]byte, error) {

	requestBody := map[string]interface{}{
		"search_type": queryOptions.Type,
		"sort_fields": queryOptions.SortFields,
		"query": map[string]interface{}{
			"term":  queryOptions.SearchTerm,
			"field": queryOptions.Field,
		},
		"from":        queryOptions.From,
		"max_results": queryOptions.MaxResults,
		"_source":     queryOptions.Source,
	}

	jsonByte, err := json.Marshal(requestBody)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, app.Config.Zs.URL, bytes.NewBuffer(jsonByte))
	if err != nil {
		app.Logger.Panicln(err)
		return nil, err
	}

	req.SetBasicAuth(app.Config.Zs.Admin, app.Config.Zs.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		app.Logger.Println("error here")
		app.Logger.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		app.Logger.Println(err)
		return nil, err
	}

	return body, nil
}
