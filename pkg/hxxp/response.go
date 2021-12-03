package hxxp

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespondJson(w http.ResponseWriter, code int, message string, data interface{}) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	dataByte, err := json.Marshal(Response{
		Message: message,
		Data:    data,
	})
	if err != nil {
		log.Println(err)
		return
	}

	_, _ = w.Write(dataByte)
}
