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
		err = errors.New(fmt.Sprint("Erro ao converter variável de ambiente OPENAI_CHATGPT_API_CHARACTER_LIMIT_PER_REQUEST "+
			"para inteiro: ", err))
		log.Error(err.Error())
		return "", err
	}

	return requestToChatGptWithCharacterLimitPerRequest(command, content, purpose, characterLimitPerRequest)
}

func requestToChatGptWithCharacterLimitPerRequest(command, content, purpose string, characterLimitPerRequest int) (string, error) {
	log.Info("Iniciando comunicação com ChatGPT: ", purpose)

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
		err = errors.New(fmt.Sprint("Erro ao converter variável de ambiente OPENAI_CHATGPT_API_REQUEST_LIMIT_PER_MINUTE"+
			" para inteiro: ", err))
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

		requestBody, err := json.Marshal(ChatGptRequest{
			Model: os.Getenv("OPENAI_CHATGPT_API_MODEL"),
			Messages: []ChatGptMessage{
				{
					Role:    "user",
					Content: fmt.Sprint(command, requestResult, partOfTheContent),
				},
			},
		})
		if err != nil {
			log.Error("Erro ao construir a requisição para comunicação com o ChatGPT: ", err)
			return "", nil
		}

		client := &http.Client{}
		request, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions",
			bytes.NewBuffer(requestBody))
		if err != nil {
			log.Error("Erro ao criar a requisição para o ChatGPT: ", err)
			return "", nil
		}
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))
		request.Header.Set("Content-Type", "application/json")

		response, err := client.Do(request)
		if err != nil {
			log.Error("Erro ao realizar requisição para o ChatGPT: ", err)
			return "", nil
		}
		defer closeResponseBody(response)

		if response.StatusCode != http.StatusOK {
			errorMessage := fmt.Sprintf("Erro a realizar requisição para o ChatGPT: (Status code: %s)", response.Status)
			log.Error(errorMessage)
			return "", errors.New(errorMessage)
		}

		var chatGptResponse ChatGptResponse
		err = json.NewDecoder(response.Body).Decode(&chatGptResponse)
		if err != nil {
			log.Error("Erro ao ler o corpo da resposta retornada pelo ChatGPT: ", err)
			return "", err
		}

		if len(chatGptResponse.Choices) < 1 {
			if response.StatusCode == 400 {
				log.Warnf("Erro ao realizar a requisição para o ChatGPT (%s), dividindo o conteúdo em partes "+
					"menores", purpose)

				characterLimit, err := strconv.Atoi(os.Getenv("OPENAI_CHATGPT_API_CHARACTER_LIMIT_PER_REQUEST"))
				if err != nil {
					log.Error("Erro ao converter variável de ambiente OPENAI_CHATGPT_API_CHARACTER_LIMIT_PER_REQUEST para "+
						"inteiro: ", err)
				} else if characterLimitPerRequest == characterLimit {
					time.Sleep(time.Minute) // To avoid excessive requests to ChatGPT
					return requestToChatGptWithCharacterLimitPerRequest(command, originalContent, purpose,
						characterLimitPerRequest/3)
				}
			}

			errorMessage := "não foi possível obter o resultado da solicitação ao ChatGPT"
			log.Error(errorMessage)
			return "", errors.New(errorMessage)
		}

		log.Infof("%dª comunicação realizada com sucesso para o ChatGPT: %s", index+1, purpose)
		requestResult = chatGptResponse.Choices[0].Message.Content
	}

	log.Info("Comunicação realizada com sucesso para o ChatGPT: ", purpose)
	time.Sleep(time.Minute) // To avoid excessive requests to ChatGPT

	return requestResult, nil
}
