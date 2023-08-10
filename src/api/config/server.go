package config

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/unidoc/unipdf/v3/common/license"
	"os"
	"time"
	"vnc-write-api/api/config/diconteiner"
	"vnc-write-api/api/endpoints/routes"
)

func NewServer() {
	loadEnvFile()

	go getBackgroundData()

	app := routes.LoadRoutes()
	address := fmt.Sprintf("%s:%s", os.Getenv("SERVER_ADDRESS"), os.Getenv("SERVER_PORT"))
	app.Logger.Fatal(app.Start(address))
}

func getBackgroundData() {
	err := license.SetMeteredKey(os.Getenv("UNI_CLOUD_KEY"))
	if err != nil {
		log.Error("Erro ao validar chave para manipulação de PDF: ", err.Error())
		return
	}

	backgroundDataService := diconteiner.GetBackgroundDataService()
	backgroundDataService.RegisterNewPropositions()

	for range time.NewTicker(time.Hour).C {
		backgroundDataService.RegisterNewPropositions()
	}
}
