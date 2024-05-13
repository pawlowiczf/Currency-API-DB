package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func fromJSON[T any](res *http.Response, target *T) error {
	//
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)

	err := json.Unmarshal(buf.Bytes(), target)

	if err != nil {
		log.Print(err)
		log.Fatalf("Error unmarshal")
	}
	return err
}

func MakeHTTPGetRequest[T any](
	rawUrl string,
	queryParameters url.Values,
	target *T) error {

	client := http.Client{}

	request, err := http.NewRequest("GET", rawUrl, nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	return fromJSON(response, target)
}
