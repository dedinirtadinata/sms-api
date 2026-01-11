package main

import (
	"log"
	"net/http"
)

func main() {
	InitDB()

	http.HandleFunc("/send-sms", SendSMSHandler)

	log.Println("SMS API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
