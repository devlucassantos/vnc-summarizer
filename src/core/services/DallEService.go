package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
)

type DallERequest struct {
	Model          string `json:"model"`
	NumberOfImages int    `json:"n"`
	Size           string `json:"size"`
	Prompt         string `json:"prompt"`
}

type DallEResponse struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}

func requestToDallE(prompt, purpose string) (string, error) {
	log.Info("Iniciando comunicação com DALL·E: ", purpose)

	requestBody, err := json.Marshal(DallERequest{
		Model:          os.Getenv("OPENAI_DALLE_API_MODEL"),
		NumberOfImages: 1,
		Size:           "1024x1024",
		Prompt:         prompt,
	})
	if err != nil {
		log.Error("Erro ao construir a requisição para comunicação com o DALL·E: ", err)
		return "", nil
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations",
		bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Erro ao criar a requisição para o DALL·E: ", err)
		return "", nil
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Error("Erro ao realizar requisição para o DALL·E: ", err)
		return "", nil
	}
	defer closeResponseBody(response)

	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Erro a realizar requisição para o DALL·E: (Status code: %s)", response.Status)
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	var dallEResponse DallEResponse
	err = json.NewDecoder(response.Body).Decode(&dallEResponse)
	if err != nil {
		log.Error("Erro ao ler o corpo da resposta retornada pelo DALL·E: ", err)
		return "", err
	}

	if len(dallEResponse.Data) < 1 {
		errorMessage := "não foi possível obter o resultado da solicitação ao DALL·E"
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	log.Info("Comunicação realizada com sucesso para o DALL·E: ", purpose)

	return dallEResponse.Data[0].Url, nil
}
