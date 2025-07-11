package dicontainer

import (
	"vnc-summarizer/adapters/apis/chatgpt"
	interfaces "vnc-summarizer/core/interfaces/chatgpt"
)

func GetChatGptApi() interfaces.ChatGpt {
	return chatgpt.NewChatGptApi()
}
