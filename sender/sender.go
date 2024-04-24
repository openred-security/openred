package sender

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	opensearch "github.com/opensearch-project/opensearch-go"
	opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

func New() *opensearch.Client {

	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{"http://localhost:9200"},
	})

	if err != nil {
		fmt.Println("failed to connect", err)
		os.Exit(1)
	}

	return client
}

func CreateIndex(client *opensearch.Client) {

	settings := strings.NewReader(`{
		'settings': {
			'index': {
				'number_of_shards': 1,
				'number_of_replicas': 0
				}
			}
		}`)

	res := opensearchapi.IndicesCreateRequest{
		Index: "go-test-index1",
		Body:  settings,
	}
	fmt.Println(res)

}
func Send(client *opensearch.Client, message string) {
	fmt.Println(message)
	document := strings.NewReader(`{
		"title": "Moneyball",
		"director": "Bennett Miller",
		"year": "2011"
	}`)

	req := opensearchapi.IndexRequest{
		Index: "go-test-index1",
		Body:  document,
	}
	insertResponse, err := req.Do(context.Background(), client)

	if err != nil {
		fmt.Println("failed to perform bulk operations", err)
	}
	fmt.Println("Performing bulk operations")
	fmt.Println(insertResponse)

}
