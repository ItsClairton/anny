package i18n

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/logger"
)

var languageMap = map[string]*Locale{}

func Load(dir string) error {

	files, err := os.ReadDir(dir)

	if err != nil {
		return err
	}

	for _, file := range files {

		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {

			buff, err := os.ReadFile(utils.Fmt("%s/%s", dir, file.Name()))
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
			logger.Debug(utils.Fmt("A Linguagem %s foi carregada com sucesso, Yeah.", info.Name))
		}

	}

	if languageMap[os.Getenv("DEFAULT_LOCALE")] == nil {
		return errors.New("invalid default locale in env path")
	}

	return nil
}

func GetDefaultLocale() *Locale {
	return languageMap[os.Getenv("DEFAULT_LOCALE")]
}

func GetLocale(id string) *Locale {
	locale, exist := languageMap[id]

	defaultLc := os.Getenv("DEFAULT_LOCALE")
	if !exist && id != defaultLc {
		logger.Warn("Não foi possível encontrar a linguagem %s, alterando para a linguagem principal.", id)
		locale = GetDefaultLocale()
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

	response, err := utils.GetFromWeb("https://translate.googleapis.com/translate_a/single?client=gtx&sl=" + url.QueryEscape(from) + "&tl=" + url.QueryEscape(to) + "&dt=t&q=" + url.QueryEscape(source))
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
			text += utils.Fmt("%v", translated)
			break
		}
	}

	return text, nil
}
