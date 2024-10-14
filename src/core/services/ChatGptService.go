package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"strconv"
	"time"
	"vnc-summarizer/core/services/utils/converters"
	"vnc-summarizer/core/services/utils/requests"
)

type ChatGptMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGptRequest struct {
	Model    string           `json:"model"`
	Messages []ChatGptMessage `json:"messages"`
}

type ChatGptResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func requestToChatGpt(command, content, purpose string) (string, error) {
	characterLimitPerRequest, err := strconv.Atoi(os.Getenv("OPENAI_CHATGPT_API_CHARACTER_LIMIT_PER_REQUEST"))
	if err != nil {
		err = errors.New(fmt.Sprint("Error converting environment variable OPENAI_CHATGPT_API_CHARACTER_LIMIT_PER_REQUEST "+
			"to integer: ", err.Error()))
		log.Error(err.Error())
		return "", err
	}

	return requestToChatGptWithCharacterLimitPerRequest(command, content, purpose, characterLimitPerRequest)
}

func requestToChatGptWithCharacterLimitPerRequest(command, content, purpose string, characterLimitPerRequest int) (string, error) {
	log.Info("Starting communication with ChatGPT: ", purpose)

	var partsOfTheContent []string
	originalContent := content
	for len(content) > 0 {
		if len(content) <= characterLimitPerRequest {
			partsOfTheContent = append(partsOfTheContent, content)
			break
		} else {
			partOfTheContent := content[:characterLimitPerRequest]
			content = content[characterLimitPerRequest:]
			partsOfTheContent = append(partsOfTheContent, partOfTheContent)
		}
	}

	requestLimitPerMinute, err := strconv.Atoi(os.Getenv("OPENAI_CHATGPT_API_REQUEST_LIMIT_PER_MINUTE"))
	if err != nil {
		err = errors.New(fmt.Sprint("Error converting environment variable OPENAI_CHATGPT_API_REQUEST_LIMIT_PER_MINUTE "+
			"to integer: ", err.Error()))
		log.Error(err.Error())
		return "", err
	}

	var requestResult string
	var requestsPerMinute int
	for index, partOfTheContent := range partsOfTheContent {
		requestsPerMinute++
		if requestsPerMinute > requestLimitPerMinute {
			time.Sleep(time.Minute) // To avoid excessive requests to ChatGPT
			requestsPerMinute = 1
		}

		body := ChatGptRequest{
			Model: os.Getenv("OPENAI_CHATGPT_API_MODEL"),
			Messages: []ChatGptMessage{
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

		request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions",
			bytes.NewBuffer(requestBody))
		if err != nil {
			log.Error("Error building the request for communication with ChatGPT: ", err.Error())
			return "", nil
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
		request.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Error("Error making request to ChatGPT: ", err.Error())
			return "", nil
		}
		defer requests.CloseResponseBody(request, response)

		if response.StatusCode != http.StatusOK {
			errorMessage := fmt.Sprintf("Error making request to ChatGPT: [Status code: %s]", response.Status)
			log.Error(errorMessage)
			return "", errors.New(errorMessage)
		}

		var chatGptResponse ChatGptResponse
		err = json.NewDecoder(response.Body).Decode(&chatGptResponse)
		if err != nil {
			log.Error("Error reading the response body returned by ChatGPT: ", err.Error())
			return "", err
		}

		if len(chatGptResponse.Choices) < 1 {
			if response.StatusCode == 400 {
				log.Warnf("Error making request to ChatGPT (%s): The content will be divided into smaller parts",
					purpose)

				characterLimit, err := strconv.Atoi(os.Getenv("OPENAI_CHATGPT_API_CHARACTER_LIMIT_PER_REQUEST"))
				if err != nil {
					log.Error("Error converting environment variable OPENAI_CHATGPT_API_CHARACTER_LIMIT_PER_REQUEST "+
						"to integer: ", err.Error())
				} else if characterLimitPerRequest == characterLimit {
					time.Sleep(time.Minute) // To avoid excessive requests to ChatGPT
					return requestToChatGptWithCharacterLimitPerRequest(command, originalContent, purpose,
						characterLimitPerRequest/3)
				}
			}

			errorMessage := "could not get the result of the request to ChatGPT"
			log.Error(errorMessage)
			return "", errors.New(errorMessage)
		}

		log.Infof("%dth successful communication with ChatGPT: %s", index+1, purpose)
		requestResult = chatGptResponse.Choices[0].Message.Content
	}

	log.Info("Successful communication with ChatGPT: ", purpose)
	time.Sleep(time.Minute) // To avoid excessive requests to ChatGPT

	return requestResult, nil
}
