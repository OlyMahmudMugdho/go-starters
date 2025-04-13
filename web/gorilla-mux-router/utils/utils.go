package utils

import (
	"encoding/json"
	"log"
)

func Json(data interface{}) ([]byte, error) {
	content, err := json.Marshal(data)
	if err != nil {
		log.Printf("error while marshalling : %v", err)
		return nil, err
	}
	return content, nil
}
