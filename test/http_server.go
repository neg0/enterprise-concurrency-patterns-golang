package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 5)
		_, _ = fmt.Fprintf(w, "Healthy Test endpoint")
	})
	http.HandleFunc("/oauth/account", func (w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 5)
		_, _ = fmt.Fprintf(w, "Account Authorisation")
	})

	http.HandleFunc("/kms", func (w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 5)
		_, _ = fmt.Fprintf(w, "Decrypted private key response result")
	})

	_ = http.ListenAndServe(":8091", nil)
}
