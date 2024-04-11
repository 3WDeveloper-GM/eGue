package app

import "flag"

type zsAPI struct {
	apiUrl      string
	admin       string
	apiPassword string
}

type ZsConfigration struct {
	api *zsAPI
}

func NewZsConfiguration() *ZsConfigration {

	api := &zsAPI{}
	flag.StringVar(&api.apiUrl, "url", "http://localhost:4080/api/", "url for the api bulk index feature")
	flag.StringVar(&api.admin, "admin", "admin", "zincsearch database administrator name")
	flag.StringVar(&api.apiPassword, "password", "Complexpass#123", "zincsearch database administrator password")

	return &ZsConfigration{
		api: api,
	}
}

func (zscfg *ZsConfigration) GetDBAdmin() string {
	return zscfg.api.admin
}

func (zscfg *ZsConfigration) GetDBPassword() string {
	return zscfg.api.apiPassword
}

func (zscfg *ZsConfigration) GetDBURL() string {
	return zscfg.api.apiUrl
}
