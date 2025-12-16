package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/search_worker/models"
)

func CreateOrUpdateTransaction(es *elasticsearch.Client, tx models.Transaction) error {
	data, _ := json.Marshal(tx)
	res, err := es.Index(
		"transactions",
		bytes.NewReader(data),
		es.Index.WithDocumentID(strconv.Itoa(tx.ID)),
		es.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Println("Indexed transaction:", tx.ID)
	return nil
}

func DeleteTransaction(es *elasticsearch.Client, tx models.Transaction) error {
	tx.Status = "reverted"

	data, _ := json.Marshal(tx)
	res, err := es.Index(
		"transactions",
		bytes.NewReader(data),
		es.Index.WithDocumentID(strconv.Itoa(tx.ID)),
		es.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Println("Reverted transaction:", tx.ID)
	return nil
}
