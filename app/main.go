package main

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"log"
	"wordsapi/app/server"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	es, err := getESClient()
	if err != nil {
		log.Fatal(err)
	}

	app := server.NewServer(es)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func getESClient() (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("http://%s:%s", viper.GetString("elasticsearch.host"), viper.GetString("elasticsearch.port"))),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	return client, err
}