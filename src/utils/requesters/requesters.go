package requesters

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"net/http"
)

func GetRequest(url string) (*http.Response, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Errorf("Error making request to %s: %s", url, err.Error())
		return nil, err
	}

	return response, err
}

func DecodeResponseBody(response *http.Response) (map[string]interface{}, error) {
	content := make(map[string]interface{})
	err := json.NewDecoder(response.Body).Decode(&content)
	if err != nil {
		log.Error("Error decoding response body: ", err.Error())
		return nil, err
	}

	return content, nil
}

func CloseResponseBody(request *http.Request, response *http.Response) {
	err := response.Body.Close()
	if err != nil {
		log.Warn("Error closing the request response body: ", request)
	}
}
