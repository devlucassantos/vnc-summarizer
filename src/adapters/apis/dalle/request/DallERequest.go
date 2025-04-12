package request

type DallERequest struct {
	Model          string `json:"model"`
	NumberOfImages int    `json:"n"`
	Size           string `json:"size"`
	Prompt         string `json:"prompt"`
}
