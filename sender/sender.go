package sender

import (
    "context"
    "crypto/tls"
	"net/http"
    "strings"
	"fmt"
	"os"
	opensearch "github.com/opensearch-project/opensearch-go"
    opensearchapi "github.com/opensearch-project/opensearch-go/opensearchapi"
)

func New(){

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

funct CreateIndex(client Client){

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
func Sender(client Client, message string) {
	fmt.Println(message)
	document := strings.NewReader(`{
		"title": "Moneyball",
		"director": "Bennett Miller",
		"year": "2011"
	}`)
	
	docId := "1"
	req := opensearchapi.IndexRequest{
		Index:      "go-test-index1",
		DocumentID: docId,
		Body:       document,
	}
	insertResponse, err := req.Do(context.Background(), client)

	if err != nil {
        fmt.Println("failed to perform bulk operations", err)
        os.Exit(1)
    }
    fmt.Println("Performing bulk operations")
    fmt.Println(insertResponse)

}
