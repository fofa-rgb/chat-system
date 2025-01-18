package database

import (
	"context"
	"fmt"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DATABASE *sqlx.DB
var ESClient *elasticsearch.Client

const SearchIndex = "messages"

func InitDB() error {
	if DATABASE != nil {
		return nil
	}
	userName := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	passWord := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", userName, passWord, host, port, name)
	database, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return err
	}
	DATABASE = database
	return nil
}

func ESClientConnection() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
		},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("Error creating Elasticsearch client: %s", err))
	}
	ESClient = client
}

func ESCreateIndexIfNotExist() {
	_, err := esapi.IndicesExistsRequest{
		Index: []string{SearchIndex},
	}.Do(context.Background(), ESClient)

	if err != nil {
		ESClient.Indices.Create(SearchIndex)
	}
}
