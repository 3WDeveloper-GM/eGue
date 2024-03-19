package handlers

import (
	"fmt"
	"net/http"

	"github.com/3WDeveloper-GM/eGue/cmd/config"
	"github.com/3WDeveloper-GM/eGue/cmd/utils"
	zincsearch "github.com/3WDeveloper-GM/eGue/cmd/zincSearch"
)

func GetHealthCheckHandler(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		response := &struct {
			Status      string `json:"status"`
			Environment string `json:"env"`
			Port        int    `json:"port"`
		}{
			Status:      "available",
			Environment: app.Config.Env,
			Port:        app.Config.Port,
		}

		formattedResponse := utils.Envelope{"response": response}

		err := utils.WriteJSON(w, http.StatusOK, formattedResponse, nil)
		if err != nil {
			app.Logger.Println(err)
			http.Error(w, "server encountered a problem", http.StatusInternalServerError)
		}

	}
}

func PostSearchEmailMatchphrase(app *config.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input = &zincsearch.Search{}

		err := utils.ReadJSON(w, r, input)
		if err != nil {
			app.Logger.Println(err)
			return
		}

		// err = utils.WriteJSON(w, http.StatusAccepted, input, nil)
		// if err != nil {
		// 	app.Logger.Println(err)
		// 	return
		// }

		result, err := zincsearch.MatchPhrase(app, input)
		if err != nil {
			app.Logger.Println(err)
			return
		}

		fmt.Fprintln(w, string(result))

	}
}
