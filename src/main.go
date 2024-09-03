package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/unidoc/unipdf/v3/common/license"
	"os"
	"time"
	"vnc-summarizer/config/diconteiner"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatal("Arquivo de variáveis de ambiente não encontrado: ", err)
	}

	err = license.SetMeteredKey(os.Getenv("UNICLOUD_KEY"))
	if err != nil {
		log.Error("Erro ao validar chave para manipulação de PDF: ", err.Error())
		return
	}

	backgroundDataService := diconteiner.GetBackgroundDataService()
	backgroundDataService.RegisterNewPropositions()

	for range time.NewTicker(time.Hour).C {
		backgroundDataService.RegisterNewPropositions()
		timeNow := time.Now()
		if timeNow.Hour() >= 18 {
			backgroundDataService.RegisterNewNewsletter(timeNow)
		} else if timeNow.Hour() < 6 {
			backgroundDataService.RegisterNewNewsletter(timeNow.AddDate(0, 0, -1))
		}
	}
}
