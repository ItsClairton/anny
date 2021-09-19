package providers

import (
	"os"

	"github.com/ItsClairton/Anny/utils"
	"github.com/buger/jsonparser"
)

func GetRandomCat(gif bool) (string, error) {
	if os.Getenv("THECATAPI_KEY") == "NONE" {
		return getFromNekosLife("meow")
	}
	if gif {
		return getFromTheCat(true)
	}

	if utils.RandomBool() {
		return getFromNekosLife("meow")
	} else {
		return getFromTheCat(false)
	}
}

// https://nekos.life
func getFromNekosLife(typ string) (string, error) {
	json, err := utils.GetFromWeb(utils.Fmt("https://nekos.life/api/v2/img/%s", typ))
	if err != nil {
		return "", err
	}

	return jsonparser.GetString(json, "url")
}

// https://thecatapi.com
func getFromTheCat(gif bool) (string, error) {
	json, err := utils.GetFromWeb(utils.Fmt("https://api.thecatapi.com/v1/images/search?api_key=%s&format=json%s", os.Getenv("THECATAPI_KEY"), utils.Is(gif, "&mime_types=gif", "")))
	if err != nil {
		return "", err
	}

	return jsonparser.GetString(json, "[0]", "url")
}
