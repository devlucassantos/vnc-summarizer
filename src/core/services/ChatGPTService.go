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

func requestToChatGpt(content, purpose string) (string, error) {
	log.Info("Iniciando comunicação com ChatGPT: ", purpose)

	requestBody, err := json.Marshal(ChatGptRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatGptMessage{
			{
				Role:    "user",
				Content: content,
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
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CHAT_GPT_KEY")))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Error("Erro ao realizar requisição para o ChatGPT: ", err)
		return "", nil
	}
	defer closeResponseBody(response)

	var chatGptResponse ChatGptResponse
	err = json.NewDecoder(response.Body).Decode(&chatGptResponse)
	if err != nil {
		log.Error("Erro ao ler o corpo da resposta retornada pelo ChatGPT: ", err)
		return "", err
	}

	if len(chatGptResponse.Choices) < 1 {
		log.Error("Não foi possível obter o resultado da solicitação ao ChatGPT")
		return "", errors.New("não foi possível obter o resultado da solicitação ao ChatGPT")
	}

	log.Info("Comunicação realizada com sucesso para o ChatGPT: ", purpose)
	return chatGptResponse.Choices[0].Message.Content, nil
}
