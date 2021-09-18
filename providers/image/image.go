package image

import (
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"

	"github.com/ItsClairton/Anny/providers/anilist"
	"github.com/ItsClairton/Anny/utils"
	"github.com/buger/jsonparser"
)

var (
	client = &http.Client{}
)

type TraceEntry struct {
	Title        *anilist.MediaTitle
	Adult        bool
	Episode      int64
	From, To     float64
	Video, Image string
}

func GetFromTrace(mediaUrl string) (*TraceEntry, error) {

	response, err := utils.GetFromWeb(utils.Fmt("https://api.trace.moe/search?cutBorders&anilistInfo&url=%s", url.QueryEscape(mediaUrl)))

	if err != nil {
		return nil, err
	}

	traceErr, err := jsonparser.GetString(response, "error")

	if err != nil {
		return nil, err
	}

	if traceErr != "" {
		return nil, errors.New(utils.Fmt("trace.moe error: %s", traceErr))
	}

	episode, _ := jsonparser.GetInt(response, "result", "[0]", "episode")
	jpName, _ := jsonparser.GetString(response, "result", "[0]", "anilist", "title", "romaji")
	enName, _ := jsonparser.GetString(response, "result", "[0]", "anilist", "title", "english")
	from, _ := jsonparser.GetFloat(response, "result", "[0]", "from")
	to, _ := jsonparser.GetFloat(response, "result", "[0]", "to")
	video, _ := jsonparser.GetString(response, "result", "[0]", "video")
	image, _ := jsonparser.GetString(response, "result", "[0]", "image")
	adult, _ := jsonparser.GetBoolean(response, "result", "[0]", "anilist", "isAdult")

	return &TraceEntry{
		Title: &anilist.MediaTitle{
			JP: jpName,
			EN: enName,
		},
		Adult:   adult,
		Episode: episode,
		From:    from,
		To:      to,
		Video:   video,
		Image:   image,
	}, nil

}

func GetFromNekos(nType string) (string, error) {
	json, err := utils.GetFromWeb(utils.Fmt("https://nekos.life/api/v2/img/%s", nType))

	if err != nil {
		return "", err
	}

	return jsonparser.GetString(json, "url")
}

func GetFromNekoBot(nType string) (string, error) {
	json, err := utils.GetFromWeb(utils.Fmt("https://nekobot.xyz/api/image?type=%s", nType))

	if err != nil {
		return "", err
	}

	return jsonparser.GetString(json, "message")
}

func GetRandomCat(gif bool) (string, error) {
	if !gif && rand.Float32() < 0.5 {
		return GetFromNekos("meow")
	}

	req, _ := http.NewRequest("GET", utils.Fmt("https://api.thecatapi.com/v1/images/search?format=json%s", utils.Is(gif, "&mime_types=gif", "")), nil)
	req.Header.Set("x-api-key", os.Getenv("CATAPI_KEY"))
	res, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return jsonparser.GetString(body, "[0]", "url")
}
