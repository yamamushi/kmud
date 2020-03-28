package main

import (
	"github.com/yamamushi/kmud-2020/database"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/yamamushi/kmud-2020/config"
	"github.com/yamamushi/kmud-2020/utils"
)

func main() {

	log.Println("Checking config")
	conf, err := config.GetConfig("accountmanager.conf")
	if err != nil {
		utils.HandleError(err)
	}

	db := database.NewDatabaseHandler(conf)

	log.Println("Creating endpoint handlers")
	svc := accountManagerService{}

	// Auth
	authHandler := httptransport.NewServer(
		makeAuthEndpoint(svc, conf, db),
		decodeAuthRequest,
		encodeResponse,
	)

	// Account Info
	accountInfoHandler := httptransport.NewServer(
		makeAccountInfoEndpoint(svc, conf, db),
		decodeAccountInfoRequest,
		encodeResponse,
	)

	// Register Account
	accountRegistrationHandler := httptransport.NewServer(
		makeAccountRegistrationEndpoint(svc, conf, db),
		decodeAccountRegistrationRequest,
		encodeResponse,
	)

	// Account Search
	searchHandler := httptransport.NewServer(
		makeSearchEndpoint(svc, conf, db),
		decodeSearchRequest,
		encodeResponse,
	)

	// Modify Account
	modifyHandler := httptransport.NewServer(
		makeModifyEndpoint(svc, conf, db),
		decodeModifyRequest,
		encodeResponse,
	)
	/*
		- ModifyPermissionsGroups
	*/

	log.Println("Registering endpoint handlers")
	http.Handle("/auth", authHandler)
	http.Handle("/accountinfo", accountInfoHandler)
	http.Handle("/modify", modifyHandler)
	http.Handle("/register", accountRegistrationHandler)
	http.Handle("/search", searchHandler)

	log.Println("Listening for connections...")
	err = http.ListenAndServe(conf.Server.Interface+":"+conf.Server.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
