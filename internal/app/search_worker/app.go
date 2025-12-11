package searchworker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/segmentio/kafka-go"

	config "github.com/go-park-mail-ru/2025_2_VKarmane/cmd/api/app"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/search_worker/handlers"
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
			log.Fatal("Kafka read error:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		var wrapper models.KafkaMessageWrapper
		if err := json.Unmarshal(m.Value, &wrapper); err != nil {
			log.Fatal("JSON unmarshal error:", err)
			continue
		}

		switch wrapper.Type {
		case models.TRANSACTIONS:
			var tx models.Transaction
			if err := (&tx).UnmarshalJSON(wrapper.Payload); err != nil {
				log.Fatal("easyjson unmarshal error:", err)
				continue
			}
			switch tx.Action {
			case "create", "update":
				if err := handlers.CreateOrUpdateTransaction(es, tx); err != nil {
					log.Fatal("Elasticsearch index error:", err)
				}
			case "delete":
				if err := handlers.DeleteTransaction(es, tx); err != nil {
					log.Fatal("Elasticsearch delete error:", err)
				}
			default:
				log.Fatal("Unknown action:", tx.Action)
			}
		case models.CATEGORIES:
			var ctg models.Category
			if err := (&ctg).UnmarshalJSON(wrapper.Payload); err != nil {
				log.Fatal("easyjson unmarshal error:", err)
				continue
			}
			switch ctg.Action {
			case "update":
				if err := handlers.UpdateCategoryInTransaction(es, ctg); err != nil {
					log.Fatal("Elasticsearch index error:", err)
				}
			default:
				log.Fatal("Unknown action:", ctg.Action)
			}
		}
	}

}
