package utils

import (
	"io"
	"net/http"
)

func GetFromWebString(url string) (string, error) {
	body, err := GetFromWeb(url)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func GetFromWeb(url string) ([]byte, error) {
	var body []byte
	res, err := http.Get(url)

	if err != nil {
		return body, err
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
