package main

import (
	"flag"
	"log"
	"os"

	"github.com/3WDeveloper-GM/eGue/cmd/config"
	"github.com/3WDeveloper-GM/eGue/cmd/indexer"
)

func main() {

	conf := &config.Config{
		Zs: config.ZSConfig{},
	}

	flag.IntVar(&conf.Port, "port", 4040, "defines the port value")
	flag.StringVar(&conf.Env, "env", "development", "environment(development|production|staging)")
	flag.StringVar(&conf.Zs.URL, "searchURL", "http://localhost:4080/api/mails/_search", "defines the search URL for ZS")
	flag.StringVar(&conf.Zs.Admin, "admin", "admin", "defines the user to access ZincSearch")
	flag.StringVar(&conf.Zs.Password, "pass", "Complexpass#123", "define the password for the user")
	flag.StringVar(&conf.IndexPath, "ZincApi", "http://localhost:4080/index", "defines the amount of indexes")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &config.Application{
		Config: *conf,
		Logger: logger,
	}

	indexer.StartIndexing()

	err := startServer(app)
	if err != nil {
		log.Fatal(err)
		return
	}
}
