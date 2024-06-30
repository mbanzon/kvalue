package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mbanzon/kvalue/internal/storage"
)

type Server struct {
	start func() error
}

type ConfigFunc func(*Server)

func New(storage *storage.Storage, port int, writeKey string) *Server {
	getHandler := createGetValueHandler(storage)
	saveHandler := emptySaveValueHandler

	if writeKey != "" {
		saveHandler = protectSaveValueHandler(writeKey, createSaveValueHandler(storage))
	}

	s := &Server{
		start: createServerStarter(getHandler, saveHandler, port),
	}

	return s
}

func emptySaveValueHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func protectSaveValueHandler(writeKey string, f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != writeKey {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		f(w, r)
	})
}

func createSaveValueHandler(s *storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "error decoding body", http.StatusBadRequest)
			return
		}

		err = s.Store(key, tmp)
		if err != nil {
			http.Error(w, "error storing value", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func createGetValueHandler(s *storage.Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key is required", http.StatusBadRequest)
			return
		}

		value, err := s.Get(key)
		if err != nil {
			if err == storage.ErrNotFound {
				http.Error(w, "not found", http.StatusNotFound)
				return
			} else {
				http.Error(w, "error getting value", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(value)
	})
}

func createMainHandler(getHandler, saveHandler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		if r.Method == http.MethodGet {
			getHandler.ServeHTTP(w, r)
			return
		} else if r.Method == http.MethodPost {
			saveHandler.ServeHTTP(w, r)
			return
		}

		http.NotFound(w, r)
	})
}

func createServerStarter(getHandler, saveHandler http.HandlerFunc, port int) func() error {
	return func() error {
		return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), createMainHandler(getHandler, saveHandler))
	}
}

func (s *Server) Start() error {
	return s.start()
}
