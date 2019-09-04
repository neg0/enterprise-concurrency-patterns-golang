package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 5)
		_, _ = fmt.Fprintf(w, "Golang Concurrency Ninja!")
	})
	http.HandleFunc("/oauth/account", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 5)
		_, _ = fmt.Fprintf(w, "Account Authorisation")
	})
	http.HandleFunc("/kms", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond * 5)
		_, _ = fmt.Fprintf(w, "Decrypted private key response result")
	})

	http.HandleFunc("/enrich/merchant", func(writer http.ResponseWriter, request *http.Request) {
		request.Header.Set("content-type", "text/plain")
		writer.WriteHeader(http.StatusOK)

		if request.Method != "POST" {
			writer.WriteHeader(http.StatusNotFound)
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			_, _ = fmt.Fprint(writer)
			return
		}

		randMerchants := []string{
			"Netflix",
			"Mark & Spencer",
			"Amazon",
		}
		_, _ = fmt.Fprint(writer, randMerchants[rand.Intn(len(randMerchants))])
		return
	})

	http.HandleFunc("/enrich/category", func(writer http.ResponseWriter, request *http.Request) {
		request.Header.Set("content-type", "text/plain")
		writer.WriteHeader(http.StatusOK)

		if request.Method != "POST" {
			writer.WriteHeader(http.StatusNotFound)
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)

			_, _ = fmt.Fprint(writer)
			return
		}

		randCats := []string{
			"Entertainment",
			"Leisure",
			"Shopping",
		}
		_, _ = fmt.Fprint(writer, randCats[rand.Intn(len(randCats))])
		return
	})

	_ = http.ListenAndServe(":8091", nil)
}
