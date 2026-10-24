package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	mailer, err := NewInfraiEmail()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var signup Signup
		if err := json.NewDecoder(r.Body).Decode(&signup); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		result, err := SendVerification(context.Background(), mailer, signup, "http://localhost:8080")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
