package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"net/http"
	"os"
	"vnc-summarizer/core/services/utils/converters"
	"vnc-summarizer/core/services/utils/requesters"
)

func requestToVncPdfContentExtractorApi(pdfUrl string) (string, error) {
	log.Info("Starting communication with VNC PDF Content Extractor API: Extraction of PDF content available at ",
		pdfUrl)

	pdfContentExtractorAddress := fmt.Sprintf("%s/api/v1/extract-content", os.Getenv("VNC_PDF_CONTENT_EXTRACTOR_API_ADDRESS"))
	body := map[string]string{
		"pdf_url": pdfUrl,
	}
	requestBody, err := converters.ToJson(body)
	if err != nil {
		log.Error("converters.ToJson(): ", err.Error())
		return "", err
	}

	request, err := http.NewRequest("POST", pdfContentExtractorAddress, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Error("Error building the request for communication with VNC PDF Content Extractor API: ", err.Error())
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error("Error making request to VNC PDF Content Extractor API: ", err.Error())
		return "", err
	}
	defer requesters.CloseResponseBody(request, response)

	if response.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(response.Body)
		if err != nil {
			log.Error("Error interpreting VNC PDF Content Extractor API response: ", err.Error())
			return "", err
		}
		responseBodyAsString := string(responseBody)

		errorMessage := fmt.Sprintf("Error making request to VNC PDF Content Extractor API: [Status: %s; "+
			"Body: %s]", response.Status, responseBodyAsString)
		log.Error(errorMessage)
		return "", errors.New(errorMessage)
	}

	responseBody, err := requesters.DecodeResponseBody(response)
	if err != nil {
		log.Error("requests.DecodeResponseBody(): ", err.Error())
		return "", err
	}

	propositionContent := fmt.Sprint(responseBody["content"])

	log.Info("Successful communication with VNC PDF Content Extractor API: Extraction of PDF content available at ",
		pdfUrl)
	return propositionContent, nil
}
