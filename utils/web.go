package utils

import (
	"io"
	"net/http"
)

func GetFromWeb(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return io.ReadAll(res.Body)
}

func GetFromWebString(url string) (string, error) {
	body, err := GetFromWeb(url)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
