package translate

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/ItsClairton/Anny/utils/rest"
	"github.com/ItsClairton/Anny/utils/sutils"
)

func Translate(from, to, source string) (string, error) {

	if len(source) < 1 {
		return source, errors.New("empty source")
	}

	var result []interface{}
	var text string

	response, err := rest.Get("https://translate.googleapis.com/translate_a/single?client=gtx&sl=" + from + "&tl=" + to + "&dt=t&q=" + url.QueryEscape(source))
	if err != nil {
		return source, err
	}

	err = json.Unmarshal(response, &result)
	if err != nil {
		return source, err
	}

	inner := result[0]
	for _, slice := range inner.([]interface{}) {
		for _, translated := range slice.([]interface{}) {
			text += sutils.Fmt("%v", translated)
			break
		}
	}

	return text, nil
}
