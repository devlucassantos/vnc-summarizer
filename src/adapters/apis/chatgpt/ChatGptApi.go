package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"vnc-summarizer/adapters/apis/chatgpt/request"
	"vnc-summarizer/adapters/apis/chatgpt/response"
	"vnc-summarizer/utils/converters"
	"vnc-summarizer/utils/requesters"
)

type ChatGpt struct{}

func NewChatGptApi() *ChatGpt {
	return &ChatGpt{}
}

func (instance ChatGpt) MakeRequest(command, content, purpose string) (string, error) {
	log.Info("Starting communication with ChatGPT: ", purpose)

	characterLimitPerRequest, err := strconv.Atoi(os.Getenv("OPENAI_CHATGPT_API_CHARACTER_LIMIT_PER_REQUEST"))
	if err != nil {
		err = errors.New(fmt.Sprint("Error converting environment variable OPENAI_CHATGPT_API_CHARACTER_LIMIT_PER_REQUEST "+
			"to integer: ", err.Error()))
		log.Error(err.Error())
		return "", err
	}

	var requestResult string
	contentParts := getContentParts(content, characterLimitPerRequest)
	for index, partOfTheContent := range contentParts {
		body := request.ChatGptRequest{
			Model: os.Getenv("OPENAI_CHATGPT_API_MODEL"),
			Messages: []request.ChatGptMessage{
				{
					Role:    "user",
					Content: fmt.Sprint(command, requestResult, partOfTheContent),
				},
			},
		}
		requestBody, err := converters.ToJson(body)
		if err != nil {
			log.Error("converters.ToJson(): ", err.Error())
			return "", err
		}

		requestToChatGpt, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions",
			bytes.NewBuffer(requestBody))
		if err != nil {
			log.Error("Error building the request for communication with ChatGPT: ", err.Error())
			return "", nil
		}
		requestToChatGpt.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
		requestToChatGpt.Header.Set("Content-Type", "application/json")

		client := &http.Client{
			Timeout: time.Minute,
		}
		responseFromChatGpt, err := client.Do(requestToChatGpt)
		if err != nil {
			log.Error("Error making request to ChatGPT: ", err.Error())
			return "", nil
		}
		defer requesters.CloseResponseBody(requestToChatGpt, responseFromChatGpt)

		time.Sleep(time.Minute) // To avoid excessive requests to ChatGPT

		if responseFromChatGpt.StatusCode != http.StatusOK {
			responseBody, err := io.ReadAll(responseFromChatGpt.Body)
			if err != nil {
				log.Error("Error interpreting ChatGPT response: ", err.Error())
				return "", err
			}

			errorMessage := fmt.Sprintf("Error making request to ChatGPT: [Status: %s; Body: %s]",
				responseFromChatGpt.Status, string(responseBody))
			log.Error(errorMessage)
			return "", errors.New(errorMessage)
		}

		var chatGptResponse response.ChatGptResponse
		err = json.NewDecoder(responseFromChatGpt.Body).Decode(&chatGptResponse)
		if err != nil {
			log.Error("Error reading the response body returned by ChatGPT: ", err.Error())
			return "", err
		}

		if len(chatGptResponse.Choices) < 1 {
			errorMessage := fmt.Sprint("Could not get the result of the request to ChatGPT: ", purpose)
			log.Error(errorMessage)
			return "", errors.New(errorMessage)
		}

		if len(partOfTheContent) > 1 {
			log.Infof("%dth successful communication with ChatGPT: %s", index+1, purpose)
		}

		requestResult = chatGptResponse.Choices[0].Message.Content
	}

	log.Info("Successful communication with ChatGPT: ", purpose)
	return requestResult, nil
}

func getContentParts(content string, characterLimitPerRequest int) []string {
	var contentParts []string
	for len(content) > 0 {
		if len(content) <= characterLimitPerRequest {
			contentParts = append(contentParts, content)
			break
		} else {
			contentPart := content[:characterLimitPerRequest]
			content = content[characterLimitPerRequest:]
			contentParts = append(contentParts, contentPart)
		}
	}

	return contentParts
}

func (instance ChatGpt) MakeRequestToVision(imageUrl string) (string, error) {
	purpose := fmt.Sprint("Description of the image available at ", imageUrl)
	log.Info("Starting communication with ChatGPT Vision: ", purpose)

	body := request.ChatGptRequest{
		Model: os.Getenv("OPENAI_CHATGPT_API_MODEL"),
		Messages: []request.ChatGptMessage{
			{
				Role: "user",
				Content: []map[string]interface{}{
					{
						"type": "text",
						"text": "Descreva a imagem de forma clara, detalhada e acessível, priorizando informações que " +
							"transmitam o contexto, a emoção e os elementos visuais importantes. Inclua detalhes como:\n" +
							"Objetos principais e secundários;\nPessoas (aparência, ações, expressões faciais, roupas);\n" +
							"Ambiente (localização, iluminação, clima, cores dominantes);\nRelações entre os elementos " +
							"da imagem;\nQualquer texto presente na imagem.\nA descrição deve ser em texto corrido, " +
							"objetiva, incluir informações relevantes e ser fácil de entender para pessoas com " +
							"deficiência visual, evitando termos técnicos desnecessários ou vagas generalizações.",
					},
					{
						"type": "image_url",
						"image_url": map[string]interface{}{
							"url":    imageUrl,
							"detail": "high",
						},
					},
				},
			},
		},
	}
	requestBody, err := converters.ToJson(body)
	if err != nil {
		log.Error("converters.ToJson(): ", err.Error())
		return "", err
	}

	requestToChatGpt, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Error building the request for communication with ChatGPT Vision: ", err.Error())
		return "", nil
	}
	requestToChatGpt.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
	requestToChatGpt.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Minute,
	}
	responseFromChatGpt, err := client.Do(requestToChatGpt)
	if err != nil {
		log.Error("Error making request to ChatGPT Vision: ", err.Error())
		return "", nil
	}
	defer requesters.CloseResponseBody(requestToChatGpt, responseFromChatGpt)

	time.Sleep(time.Minute) // To avoid excessive requests to ChatGPT

	if responseFromChatGpt.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(responseFromChatGpt.Body)
		if err != nil {
			log.Error("Error interpreting ChatGPT Vision response: ", err.Error())
			return "", err
		}

		errorMessage := fmt.Sprintf("Error making request to ChatGPT Vision: [Status: %s; Body: %s]",
			responseFromChatGpt.Status, string(responseBody))
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	var chatGptResponse response.ChatGptResponse
	err = json.NewDecoder(responseFromChatGpt.Body).Decode(&chatGptResponse)
	if err != nil {
		log.Error("Error reading the response body returned by ChatGPT Vision: ", err.Error())
		return "", err
	}

	if len(chatGptResponse.Choices) < 1 {
		errorMessage := fmt.Sprint("Could not get the result of the request to ChatGPT Vision: ", purpose)
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	requestResult := chatGptResponse.Choices[0].Message.Content

	log.Info("Successful communication with ChatGPT Vision: ", purpose)
	return requestResult, nil
}
