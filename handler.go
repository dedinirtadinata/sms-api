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
	// Simpan tanda + jika ada
	hasPlus := len(n) > 0 && n[0] == '+'

	// Hapus semua karakter non-numerik
	re := regexp.MustCompile(`[^0-9]`)
	n = re.ReplaceAllString(n, "")

	if len(n) == 0 {
		return n
	}

	if n[0:1] == "0" {
		// Nomor lokal Indonesia (dimulai dengan 0), ubah ke format internasional
		n = "+62" + n[1:]
	} else if n[0:2] == "62" {
		// Nomor sudah dalam format internasional (dimulai dengan 62), tambahkan +
		n = "+" + n
	} else if hasPlus {
		// Jika awalnya ada +, tambahkan kembali
		n = "+" + n
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

	number := normalizeNumber(req.To)

	// Field "Text" is NULL, only "TextDecoded" is filled (UTF-8 text)
	// This matches the manual insert pattern that works with Gammu
	_, err := db.Exec(`
		INSERT INTO outbox
		("DestinationNumber", "TextDecoded", "CreatorID")
		VALUES ($1, $2, 'SYSTEM')
	`,
		number,
		req.Message,
	)

	if err != nil {
		http.Error(w, "Failed to queue SMS", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"queued"}`))
}
