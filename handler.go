package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

type SMSRequest struct {
	To       string `json:"to"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
	DLR      string `json:"dlr"`
}

func normalizeNumber(n string) string {
	re := regexp.MustCompile(`[^0-9]`)
	n = re.ReplaceAllString(n, "")

	if n[0:1] == "0" {
		n = "62" + n[1:]
	}
	return n
}

func SendSMSHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "DB not ready", 500)
		return
	}

	var req SMSRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	if req.Priority == 0 {
		req.Priority = 0
	}
	if req.DLR == "" {
		req.DLR = "default"
	}

	number := normalizeNumber(req.To)

	_, err := db.Exec(`
		INSERT INTO outbox
		("DestinationNumber","TextDecoded","CreatorID","Priority","DeliveryReport")
		VALUES ($1,$2,'API',$3,$4)
	`,
		number,
		req.Message,
		req.Priority,
		req.DLR,
	)

	if err != nil {
		http.Error(w, "Failed to queue SMS", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"queued"}`))
}
