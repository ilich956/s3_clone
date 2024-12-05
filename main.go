package main

import (
	"log"
	"net/http"
	"os"
	"triple-s/internal/config"
	"triple-s/internal/handler"
	"triple-s/pkg/csvutil"
)

func main() {
	config.ParseFlags()

	if err := os.Mkdir(*config.Dir, 0o755); err != nil {
		log.Fatal("Failed to create folder: ", err)
	}
	csvutil.CreateNewBucketCSV()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handler.HandleGetBucket)
	mux.HandleFunc("PUT /{BucketName}", handler.HandlePutBucket)
	mux.HandleFunc("DELETE /{BucketName}", handler.HandleDeleteBucket)

	mux.HandleFunc("GET /{BucketName}/{ObjectKey}", handler.HandleGetObject)
	mux.HandleFunc("PUT /{BucketName}/{ObjectKey}", handler.HandlePutObject)
	mux.HandleFunc("DELETE /{BucketName}/{ObjectKey}", handler.HandleDeleteObject)

	if err := http.ListenAndServe(":"+*config.Port, mux); err != nil {
		log.Fatal("Failed to launch server ", err)
	}
}
