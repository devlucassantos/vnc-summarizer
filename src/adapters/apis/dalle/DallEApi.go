package dalle

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
	"time"
	"vnc-summarizer/adapters/apis/dalle/request"
	"vnc-summarizer/adapters/apis/dalle/response"
	"vnc-summarizer/utils/converters"
	"vnc-summarizer/utils/requesters"
)

type DallE struct{}

func NewDallEApi() *DallE {
	return &DallE{}
}

func (instance DallE) MakeRequest(prompt, purpose string) (string, error) {
	log.Info("Starting communication with DALL·E: ", purpose)

	body := request.DallERequest{
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

	requestToDallE, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations",
		bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Error building the request for communication with DALL·E: ", err.Error())
		return "", nil
	}
	requestToDallE.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
	requestToDallE.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Minute,
	}
	responseFromDallE, err := client.Do(requestToDallE)
	if err != nil {
		log.Error("Error making request to DALL·E: ", err.Error())
		return "", nil
	}
	defer requesters.CloseResponseBody(requestToDallE, responseFromDallE)

	if responseFromDallE.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(responseFromDallE.Body)
		if err != nil {
			log.Error("Error interpreting DALL·E response: ", err.Error())
			return "", err
		}

		errorMessage := fmt.Sprintf("Error making request to DALL·E: [Status: %s; Body: %s]",
			responseFromDallE.Status, string(responseBody))
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	var dallEResponse response.DallEResponse
	err = json.NewDecoder(responseFromDallE.Body).Decode(&dallEResponse)
	if err != nil {
		log.Error("Error reading the response body returned by DALL·E: ", err.Error())
		return "", err
	}

	if len(dallEResponse.Data) < 1 {
		errorMessage := fmt.Sprint("Could not get the result of the request to DALL·E: ", purpose)
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	requestResult := dallEResponse.Data[0].Url

	log.Info("Successful communication with DALL·E: ", purpose)
	return requestResult, nil
}
