package main

import (
	"log"
	"net/http"
	"wxcloudrun-golang/service"
)

func main() {
	//if err := db.Init(); err != nil {
	//	panic(fmt.Sprintf("mysql init failed with %+v", err))
	//}

	http.HandleFunc("/api/transfer", service.TransferHandler)

	log.Fatal(http.ListenAndServe(":80", nil))
}
