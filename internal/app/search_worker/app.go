package searchworker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/segmentio/kafka-go"

	config "github.com/go-park-mail-ru/2025_2_VKarmane/cmd/api/app"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/search_worker/models"
)

func Run() error {
	config := config.LoadConfig()

	es, err := elasticsearch.NewClient(elasticsearch.Config{
    Addresses: []string{
        fmt.Sprintf("http://%s:%s", config.ElasticSearch.Host, config.ElasticSearch.Port),
    	},
	})
	if err != nil {
		log.Fatal("Failed to connect to ElasticSearch", "error", err)
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%s", config.KafkaProducerHost, config.KafkaProducerPort)},
		Topic:   "transactions",
		GroupID: "transaction-workers",
	})

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		var tx models.Transaction
		if err := json.Unmarshal(m.Value, &tx); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		data, _ := json.Marshal(tx)
		res, err := es.Index(
			"transactions",
			bytes.NewReader(data),
			es.Index.WithDocumentID(strconv.Itoa(tx.ID)),
			es.Index.WithRefresh("true"),
		)
		if err != nil {
			log.Println("Elasticsearch index error:", err)
			continue
		}
		res.Body.Close()
		log.Println("Indexed transaction:", tx.ID)
	}

}
