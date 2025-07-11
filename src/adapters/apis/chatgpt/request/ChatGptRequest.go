package request

type ChatGptRequest struct {
	Model    string           `json:"model"`
	Messages []ChatGptMessage `json:"messages"`
}

type ChatGptMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}
