package i18n

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/ItsClairton/Anny/utils/rest"
	"github.com/ItsClairton/Anny/utils/sutils"
)

var languageMap = map[string]*Locale{}

func Load(dir string) error {

	files, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, file := range files {

		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {

			buff, err := os.ReadFile(sutils.Fmt("%s/%s", dir, file.Name()))
			if err != nil {
				return err
			}

			var info *Locale

			err = json.Unmarshal(buff, &info)
			if err != nil {
				return err
			}

			info.ID = strings.TrimSuffix(file.Name(), ".json")
			info.Content = buff
			languageMap[info.ID] = info
			logger.Debug(sutils.Fmt("A Linguagem %s foi carrega com sucesso, Yeah.", info.Name))
		}

	}

	if languageMap[os.Getenv("DEFAULT_LOCALE")] == nil {
		return errors.New("invalid default locale in env path")
	}

	return nil
}

func GetLocale(id string) *Locale {
	locale, exist := languageMap[id]

	if !exist {
		logger.Warn("Não foi possível encontrar a linguagem %s, alterando para a linguagem principal.", id)
		locale = languageMap[os.Getenv("DEFAULT_LOCALE")]
	}

	return locale
}

func FromGoogle(from, to, source string) (string, error) {

	if from == to {
		return source, nil
	}

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
