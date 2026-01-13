package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
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

	if len(n) == 0 {
		return n
	}

	if n[0] == '0' {
		return "+62" + n[1:]
	}

	if len(n) >= 2 && n[0:2] == "62" {
		return "+" + n
	}

	return "+" + n
}

func sanitizeSMS(text string) string {
	var buf []rune
	for _, r := range text {
		// buang control character
		if r < 32 || r == 127 {
			continue
		}
		buf = append(buf, r)
	}
	return strings.TrimSpace(string(buf))
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

	number := normalizeNumber(req.To)
	msg := sanitizeSMS(req.Message)

	// Validate message is not empty
	if req.Message == "" {
		http.Error(w, "Message cannot be empty", 400)
		return
	}

	if len(msg) == 0 {
		http.Error(w, "Invalid SMS content", 400)
		return
	}

	// Gammu requires specific status values to process SMS
	// Try simpler format first - this matches the commented pattern that should work
	_, err := db.Exec(`
		INSERT INTO outbox
		("DestinationNumber", "TextDecoded", "CreatorID", "SendingDateTime", "Coding", "Class", "RelativeValidity")
		VALUES ($1, $2, 'SYSTEM', NOW(), 'Default_No_Compression', -1, 255)
	`, number, msg)

	if err != nil {
		http.Error(w, "Failed to queue SMS: "+err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"queued"}`))
}
