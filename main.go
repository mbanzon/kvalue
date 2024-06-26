package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	defaultPort = 8080

	writeKeyEnvironment = "KVALUE_WRITE_KEY"
)

func main() {
	port := flag.Int("port", defaultPort, "port to listen on")
	flag.Parse()

	writeKey := getEnvironment()

	values := make(map[string]json.RawMessage)

	// load values from file
	data, err := os.ReadFile("values.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &values)
	if err != nil {
		panic(err)
	}

	saveValues := func(values map[string]json.RawMessage) {
		// save values to file
		data, err := json.MarshalIndent(values, "", "\t")
		if err != nil {
			log.Println(err)
		}

		err = os.WriteFile("values.json", data, 0644)
		if err != nil {
			log.Println(err)
		}
	}

	valueCh := make(chan map[string]json.RawMessage)

	go func() {
		for v := range valueCh {
			saveValues(v)
		}
	}()

	err = http.ListenAndServe(fmt.Sprintf(":%d", *port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// check if write key is correct
			if r.Header.Get("Authorization") != writeKey {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// get key from query parameter
			key := r.URL.Query().Get("key")
			// check if key is empty
			if key == "" {
				http.Error(w, "key is required", http.StatusBadRequest)
				return
			}

			tmp := json.RawMessage{}
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&tmp)
			if err != nil {
				log.Println("error decoding body:", err)
				http.Error(w, "error decoding body", http.StatusBadRequest)
				return
			}

			values[key] = tmp
			valueCh <- values

			w.WriteHeader(http.StatusCreated)
			return
		} else if r.Method == http.MethodGet {
			key := r.URL.Query().Get("key")
			if key == "" {
				http.Error(w, "key is required", http.StatusBadRequest)
				return
			}

			value, ok := values[key]
			if !ok {
				http.Error(w, "key not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(value)
			return
		}
	}))

	if err != nil {
		panic(err)
	}
}

func getEnvironment() string {
	return os.Getenv(writeKeyEnvironment)
}
