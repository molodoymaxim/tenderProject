package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"tenderProject/backend/internal/config"
	NewRouter "tenderProject/backend/internal/http/router"
	"tenderProject/backend/internal/storage"
)

func main() {
	env := flag.String("env", "", "Specify environment (e.g. 'local')")
	flag.Parse()

	isLocal := false

	if *env == "local" {
		isLocal = true
	}

	conf := config.MustLoad(isLocal)
	storageData := storage.ConnectToStorage(conf, isLocal)
	defer storageData.Close()
	storageData.CreateTables()

	//err := storageData.NewTenderStorage()
	//if err != nil {
	//	log.Fatal("failed to create tender storage")
	//}
	//err = storageData.NewVersionStorage()
	//if err != nil {
	//	log.Fatal("failed to create version storage")
	//}
	//err = storageData.CreateBidsDB()
	//if err != nil {
	//	log.Fatal("failed to create version storage")
	//}
	//
	//err = storageData.CreateBidStory()
	//if err != nil {
	//	log.Fatal("failed to create bid story storage")
	//}

	//storageData.CreateRelation()

	router := NewRouter.NewRouter(storageData)

	addr := "0.0.0.0:8080"

	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("server started on", addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Panic("failed to start server")
	}

	log.Println("stopping server")

}
