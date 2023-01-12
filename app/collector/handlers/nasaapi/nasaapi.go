package nasaapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"url-collector/interfaces"
)

const (
	defaultNasaApiKey = "DEMO_KEY"
	nasaApiUrlF       = "https://api.nasa.gov/planetary/apod?api_key=%s&date=%s"
)

var (
	NasaApiGetError                = errors.New("Nasa API Get error")
	NasaApiReadBytesError          = errors.New("Nasa API read bytes error")
	NasaApiJsonUnmarshalError      = errors.New("Nasa API json unmarshal error")
	NasaApiErrorJsonUnmarshalError = errors.New("Nasa API error json unmarshal error")
	nasaApiKey                     = defaultNasaApiKey
)

type nasaApiResponse struct {
	Url string `json:"url"`
}

type nasaApiError struct {
	Msg string `json:"message"`
}

type nasaApiErrorResponse struct {
	Error nasaApiError `json:"error"`
}

type NasaApiUrlCollectorHandlerFactory struct {
}

type NasaApiUrlCollectorHandler struct {
	apiKey string
}

func (nauchf NasaApiUrlCollectorHandlerFactory) New() interfaces.UrlCollectorHandler {
	nasaApiUrlCollectorHandler := NasaApiUrlCollectorHandler{
		apiKey: defaultNasaApiKey,
	}

	key := os.Getenv("NASA_API_KEY")
	if key != "" {
		nasaApiUrlCollectorHandler.apiKey = key
	}

	return &nasaApiUrlCollectorHandler
}

func (nauch NasaApiUrlCollectorHandler) Handle(date string) (string, error) {
	resp, getErr := http.Get(fmt.Sprintf(nasaApiUrlF, nasaApiKey, date))
	if getErr != nil {
		return "", NasaApiGetError
	}

	body, readBytesErr := ioutil.ReadAll(resp.Body)
	if readBytesErr != nil {
		return "", NasaApiReadBytesError
	}

	var nasaApiResp nasaApiResponse
	nasaApiJsonUnmarshalErr := json.Unmarshal(body, &nasaApiResp)
	if nasaApiJsonUnmarshalErr != nil {
		return "", NasaApiJsonUnmarshalError
	}
	if nasaApiResp.Url == "" {
		var nasaApiErrorResp nasaApiErrorResponse
		nasaApiErrorJsonUnmarshalErr := json.Unmarshal(body, &nasaApiErrorResp)
		if nasaApiErrorJsonUnmarshalErr != nil {
			return "", NasaApiErrorJsonUnmarshalError
		}
		return "", errors.New(nasaApiErrorResp.Error.Msg)
	}
	return nasaApiResp.Url, nil
}
