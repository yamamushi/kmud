package main

import (
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

	log.Println("Creating endpoint handlers")
	svc := accountManagerService{}
	authHandler := httptransport.NewServer(
		makeAuthEndpoint(svc),
		decodeAuthRequest,
		encodeResponse,
	)

	accountInfoHandler := httptransport.NewServer(
		makeAccountInfoEndpoint(svc),
		decodeAccountInfoRequest,
		encodeResponse,
	)

	log.Println("Registering endpoint handlers")
	http.Handle("/auth", authHandler)
	http.Handle("/accountinfo", accountInfoHandler)

	log.Println("Listening for connections...")
	err = http.ListenAndServe(conf.Server.Interface+":"+conf.Server.Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
