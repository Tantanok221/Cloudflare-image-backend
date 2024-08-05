package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetJson[T any](result T) []byte {
	json, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Json Conversion fail", err)
	}
	return json
}

func WriteJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(GetJson(data))
}
