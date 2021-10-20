package providers

import (
	"errors"
	"net/url"
	"time"

	"github.com/ItsClairton/Anny/utils"
	"github.com/buger/jsonparser"
)

type TraceTitle struct {
	English, Japanese string
}

type TraceResult struct {
	Title        *TraceTitle
	Adult        bool
	Episode      int64
	From, To     time.Duration
	Video, Image string
}

func SearchAnimeByScene(sceneUrl string) (*TraceResult, error) {
	response, err := utils.GetFromWeb(utils.Fmt("https://api.trace.moe/search?cutBorders&anilistInfo&url=%s", url.QueryEscape(sceneUrl)))
	if err != nil {
		return nil, err
	}

	status, err := jsonparser.GetString(response, "error")
	if err != nil {
		return nil, err
	}
	if status != "" {
		return nil, errors.New(status)
	}

	episode, _ := jsonparser.GetInt(response, "result", "[0]", "episode")
	from, _ := jsonparser.GetFloat(response, "result", "[0]", "from")
	to, _ := jsonparser.GetFloat(response, "result", "[0]", "to")
	video, _ := jsonparser.GetString(response, "result", "[0]", "video")
	image, _ := jsonparser.GetString(response, "result", "[0]", "image")

	// AniList
	adult, _ := jsonparser.GetBoolean(response, "result", "[0]", "anilist", "isAdult")
	english, _ := jsonparser.GetString(response, "result", "[0]", "anilist", "title", "english")
	japanese, _ := jsonparser.GetString(response, "result", "[0]", "anilist", "title", "romaji")

	return &TraceResult{
		Title: &TraceTitle{
			English:  english,
			Japanese: japanese,
		},
		Adult:   adult,
		Episode: episode,
		From:    time.Duration(from * 1000),
		To:      time.Duration(to * 1000),
		Video:   video,
		Image:   image,
	}, nil
}
