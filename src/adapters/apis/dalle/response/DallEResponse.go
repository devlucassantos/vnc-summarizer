package response

type DallEResponse struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}
