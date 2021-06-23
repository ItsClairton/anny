package i18n

import (
	"os"
	"strings"

	"github.com/ItsClairton/Anny/logger"
	"github.com/ItsClairton/Anny/utils"
	"github.com/buger/jsonparser"
)

type Locale struct {
	ID      string
	Name    string   `json:"name"`
	Emote   string   `json:"emote"`
	Authors []string `json:"author"`
	Content []byte
}

func (lc *Locale) ToPrettyMonth(month int) string {
	return lc.GetString(utils.Fmt("months.[%d]", month-1))[0:3]
}

func (lc *Locale) ToPrettyDate(date *utils.Date) string {

	if date.Year == 0 {
		return lc.GetString("notYetReleased")
	}

	if date.Month == 0 && date.Day == 0 {
		return utils.Fmt("%d", date.Year)
	}

	if date.Day == 0 {
		return strings.TrimSpace(lc.GetString("prettyDate", lc, "", lc.ToPrettyMonth(date.Month), date.Year))
	}

	return lc.GetString("prettyDate", date.Day, lc.ToPrettyMonth(date.Month), date.Year)
}

func (lc *Locale) GetString(id string, args ...interface{}) string {

	str, err := jsonparser.GetString(lc.Content, strings.Split(id, ".")...)

	if err != nil {
		defaultLc := os.Getenv("DEFAULT_LOCALE")

		if lc.ID != defaultLc { // Fallback para a Linguagem principal.
			return GetLocale(defaultLc).GetString(id, args...)
		}
		return "N/A"
	}

	for i, content := range args {
		str = strings.ReplaceAll(str, utils.Fmt("{%d}", i), utils.Fmt("%v", content))
	}

	return str
}

func (lc *Locale) GetPrettyGenres(genres []string) []string {

	var pretty []string

	for _, genre := range genres {
		pretty = append(pretty, lc.GetPrettyGenre(strings.ToLower(genre)))
	}

	return pretty
}

func (lc *Locale) GetPrettyGenre(genre string) string {
	genreResult := lc.GetString(utils.Fmt("genres.%s", genre))

	if genreResult == "N/A" {
		logger.Warn("Não encontrei o gênero %s nos arquivos de tradução.", genre)
		return strings.Title(genre)
	}
	return genreResult
}

func (lc *Locale) GetFromArray(path string, i int) string {
	return lc.GetString(utils.Fmt("%s.[%d]", path, i))
}
