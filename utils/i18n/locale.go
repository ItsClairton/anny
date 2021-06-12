package i18n

import (
	"os"
	"strings"

	"github.com/ItsClairton/Anny/utils/logger"
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/buger/jsonparser"
)

type Locale struct {
	ID      string
	Name    string   `json:"name"`
	Emote   string   `json:"emote"`
	Authors []string `json:"author"`
	Content []byte
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
		str = strings.ReplaceAll(str, sutils.Fmt("{%d}", i), sutils.Fmt("%v", content))
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
	genreResult := lc.GetString(sutils.Fmt("genres.%s", genre))

	if genreResult == "N/A" {
		logger.Warn("Não encontrei o gênero %s nos arquivos de tradução.", genre)
		return strings.Title(genre)
	}
	return genreResult
}

func (lc *Locale) GetFromArray(path string, i int) string {
	return lc.GetString(sutils.Fmt("%s.[%d]", path, i))
}
