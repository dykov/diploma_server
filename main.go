package main

import (
	"./util"
	"log"
	"net/http"
	"time"
)

func main() {

	util.SetLogger()

	util.InitDB()

	getRouter()

	server := &http.Server{
		Handler: util.LogRequest(router),
		Addr:    ":80",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatalln(server.ListenAndServe())

}
