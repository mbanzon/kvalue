package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/mbanzon/kvalue/internal/server"
	"github.com/mbanzon/kvalue/internal/storage"
)

const (
	defaultPort = 8080

	writeKeyEnvironment = "KVALUE_WRITE_KEY"
)

func main() {
	port := flag.Int("port", defaultPort, "port to listen on")
	dataFile := flag.String("datafile", "values.json", "file to store values in")
	flag.Parse()
	writeKey := getEnvironment()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	storage, err := storage.New(*dataFile)
	if err != nil {
		log.Fatal(err)
	}

	server := server.New(storage, *port, writeKey)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func getEnvironment() string {
	return os.Getenv(writeKeyEnvironment)
}
