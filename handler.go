package main

import (
	"encoding/json"
	"net/http"
)

type SMSRequest struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func SendSMSHandler(w http.ResponseWriter, r *http.Request) {
	var req SMSRequest
	json.NewDecoder(r.Body).Decode(&req)

	_, err := db.Exec(
		"INSERT INTO outbox (DestinationNumber, TextDecoded, CreatorID) VALUES ($1,$2,'API')",
		req.To, req.Message,
	)

	if err != nil {
		http.Error(w, "Failed", 500)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))
}
