package chatgpt

type ChatGpt interface {
	MakeRequest(command, content, purpose string) (string, error)
	MakeRequestToVision(imageUrl string) (string, error)
}
