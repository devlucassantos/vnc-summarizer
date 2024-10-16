package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"vnc-summarizer/core/services/utils/converters"
	"vnc-summarizer/core/services/utils/requests"
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
	log.Info("Starting communication with DALL·E: ", purpose)

	body := DallERequest{
		Model:          os.Getenv("OPENAI_DALLE_API_MODEL"),
		NumberOfImages: 1,
		Size:           "1024x1024",
		Prompt:         prompt,
	}
	requestBody, err := converters.ToJson(body)
	if err != nil {
		log.Error("converters.ToJson(): ", err.Error())
		return "", err
	}

	request, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations",
		bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Error building the request for communication with DALL·E: ", err.Error())
		return "", nil
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error("Error making request to DALL·E: ", err.Error())
		return "", nil
	}
	defer requests.CloseResponseBody(request, response)

	if response.StatusCode != http.StatusOK {
		errorMessage := fmt.Sprintf("Error making request to DALL·E: [Status code: %s]", response.Status)
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	var dallEResponse DallEResponse
	err = json.NewDecoder(response.Body).Decode(&dallEResponse)
	if err != nil {
		log.Error("Error reading the response body returned by DALL·E: ", err.Error())
		return "", err
	}

	if len(dallEResponse.Data) < 1 {
		errorMessage := "could not get the result of the request to DALL·E"
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	log.Info("Successful communication with DALL·E: ", purpose)

	return dallEResponse.Data[0].Url, nil
}
