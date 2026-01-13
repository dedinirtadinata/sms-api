package main

import (
	"encoding/json"
	"fmt"
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
	// Ganti newline jadi spasi
	text = regexp.MustCompile(`[\r\n\t]+`).ReplaceAllString(text, " ")

	// Hapus karakter non GSM basic
	re := regexp.MustCompile(`[^\x20-\x7E]`)
	text = re.ReplaceAllString(text, "")

	// Trim spasi
	text = strings.TrimSpace(text)

	return text
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
	// Field "Text" is NULL, only "TextDecoded" is filled (UTF-8 text)
	// This matches the manual insert pattern that works with Gammu
	// _, err := db.Exec(`
	// INSERT INTO outbox
	// ("DestinationNumber", "TextDecoded", "CreatorID", "SendingDateTime", "Coding", "Class", "RelativeValidity")
	// VALUES ($1, $2, 'SYSTEM', NOW(), 'Default_No_Compression', -1, 255)
	// `,
	// 	number,
	// 	req.Message,
	// )

	_, err := db.Exec(`INSERT INTO "public"."outbox" ("SendBefore", "SendAfter", "Text", "DestinationNumber", "Coding", "UDH", "Class", "TextDecoded",  "MultiPart", "RelativeValidity", "SenderID", "SendingTimeOut", "DeliveryReport", "CreatorID", "Retries", "Priority", "Status", "StatusCode") VALUES ('23:59:59', '00:00:00', NULL,$1, 'Default_No_Compression', NULL, -1,'`+fmt.Sprintf("%s", msg)+`', 'f', -1, NULL, NOW(), 'default', 'SYSTEM', 0, 0, 'Reserved', -1);`, number)

	if err != nil {
		http.Error(w, "Failed to queue SMS", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"queued"}`))
}
