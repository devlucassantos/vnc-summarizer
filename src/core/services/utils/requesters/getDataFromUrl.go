package requesters

import (
	"github.com/labstack/gommon/log"
	"vnc-summarizer/core/services/utils/converters"
)

func getResponseFromUrl(url string) (map[string]interface{}, error) {
	response, err := GetRequest(url)
	if err != nil {
		log.Error("GetRequest(): ", err.Error())
		return nil, err
	}

	responseBody, err := DecodeResponseBody(response)
	if err != nil {
		log.Error("DecodeResponseBody(): ", err.Error())
		return nil, err
	}

	return responseBody, nil
}

func GetDataObjectFromUrl(url string) (map[string]interface{}, error) {
	responseBody, err := getResponseFromUrl(url)
	if err != nil {
		log.Error("getResponseFromUrl(): ", err.Error())
		return nil, err
	}

	resultMap, err := converters.ToMap(responseBody["dados"])
	if err != nil {
		log.Error("converters.ToMap(): ", err.Error())
		return nil, err
	}

	return resultMap, nil
}

func GetDataSliceFromUrl(url string) ([]map[string]interface{}, error) {
	responseBody, err := getResponseFromUrl(url)
	if err != nil {
		log.Error("getResponseFromUrl(): ", err.Error())
		return nil, err
	}

	resultMapSlice, err := converters.ToMapSlice(responseBody["dados"])
	if err != nil {
		log.Error("converters.ToMapSlice(): ", err.Error())
		return nil, err
	}

	return resultMapSlice, nil
}
