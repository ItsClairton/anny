package anilist

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Query struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

type Data struct {
	Media *Media `json:"Media"`
}

type Error struct {
	Message string `json:"message"`
}

type Result struct {
	Data   Data
	Errors []Error
}

func Get(query Query) (*Data, error) {

	payload, err := json.Marshal(query)

	if err != nil {
		return nil, err
	}

	raw, err := http.Post("https://graphql.anilist.co", "application/json", bytes.NewBuffer(payload))

	if err != nil {
		return nil, err
	}

	defer raw.Body.Close()

	body, err := ioutil.ReadAll(raw.Body)

	if err != nil {
		return nil, err
	}

	var result Result
	err = json.Unmarshal(body, &result)

	if err != nil {
		return nil, err
	}

	if len(result.Errors) != 0 {
		if result.Errors[0].Message == "Not Found." {
			return nil, nil
		}

		err = errors.New(result.Errors[0].Message)
		return nil, err
	}

	return &result.Data, nil
}
