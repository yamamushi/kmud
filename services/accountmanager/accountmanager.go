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
	authHandler := httptransport.NewServer(
		makeAuthEndpoint(svc, conf, db),
		decodeAuthRequest,
		encodeResponse,
	)

	accountInfoHandler := httptransport.NewServer(
		makeAccountInfoEndpoint(svc, conf, db),
		decodeAccountInfoRequest,
		encodeResponse,
	)

	accountRegistrationHandler := httptransport.NewServer(
		makeAccountRegistrationEndpoint(svc, conf, db),
		decodeAccountRegistrationRequest,
		encodeResponse,
	)

	log.Println("Registering endpoint handlers")
	http.Handle("/auth", authHandler)
	http.Handle("/accountinfo", accountInfoHandler)
	http.Handle("/register", accountRegistrationHandler)

	log.Println("Listening for connections...")
	err = http.ListenAndServe(conf.Server.Interface+":"+conf.Server.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
