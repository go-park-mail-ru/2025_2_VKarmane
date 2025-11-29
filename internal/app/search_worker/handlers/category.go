package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/search_worker/models"
)

type CategoryUpdatePayload struct {
	CategoryName         *string `json:"category_name,omitempty"`
	CategoryLogoHashedID *string `json:"category_logo_hashed_id,omitempty"`
	CategoryLogo         *string `json:"category_logo,omitempty"`
}

func UpdateCategoryInTransaction(es *elasticsearch.Client, ctg models.Category) error {
	scriptLines := []string{}
	params := map[string]interface{}{}

	if ctg.CategoryName != "" {
		scriptLines = append(scriptLines, "ctx._source.category_name = params.category_name;")
		params["category_name"] = ctg.CategoryName
	}
	if ctg.CategoryLogoHashedID != "" {
		scriptLines = append(scriptLines, "ctx._source.category_logo_hashed_id = params.category_logo_hashed_id;")
		params["category_logo_hashed_id"] = ctg.CategoryLogoHashedID
	}
	if ctg.CategoryLogo != "" {
		scriptLines = append(scriptLines, "ctx._source.category_logo = params.category_logo;")
		params["category_logo"] = ctg.CategoryLogo
	}

	if len(scriptLines) == 0 {
		return nil
	}

	script := map[string]interface{}{
		"source": strings.Join(scriptLines, " "),
		"params": params,
	}

	req := map[string]interface{}{
		"script": script,
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": ctg.CategoryID,
			},
		},
	}

	body, _ := json.Marshal(req)

	res, err := es.UpdateByQuery(
		[]string{"transactions"},
		es.UpdateByQuery.WithBody(bytes.NewReader(body)),
		es.UpdateByQuery.WithRefresh(true),
	)
	if err != nil {
		return fmt.Errorf("ES update error: %w", err)
	}
	defer res.Body.Close()

	log.Printf("Updated operations for category_id=%d\n", ctg.CategoryID)
	return nil
}