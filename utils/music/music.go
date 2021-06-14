package music

import (
	"errors"
	"net/url"

	"github.com/ItsClairton/Anny/utils/rest"
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/buger/jsonparser"
	"github.com/kkdai/youtube/v2"
)

var (
	players      = map[string]*Player{}
	client       = youtube.Client{}
	YouTubeRegex = sutils.GetRegex(`^(http(s)?:\/\/)?((w){3}.)?youtu(be|.be)?(\.com)?\/.+`)
)

func GetIDFromYouTube(content string) (string, error) {

	result, err := rest.Get(sutils.Fmt("https://youtube-scrape.herokuapp.com/api/search?q=%s", url.QueryEscape(content)))

	if err != nil {
		return "", err
	}

	id, err := jsonparser.GetString(result, "results", "[0]", "video", "id")

	if err != nil {
		return "", nil
	}

	return id, nil
}

func GetTrackFromYouTube(url string) (*Track, error) {

	info, err := client.GetVideo(url)

	if err != nil {
		return nil, err
	}

	var audioFormat *youtube.Format
	var isOpus bool

	if len(info.Formats.Type("opus")) < 1 {
		if len(info.Formats.Type("audio")) < 1 {
			return nil, errors.New("not found audio format")
		}
		audioFormat = &info.Formats.Type("audio")[0]
		isOpus = false
	} else {
		audioFormat = &info.Formats.Type("opus")[0]
		isOpus = true
	}

	resultUrl, err := client.GetStreamURL(info, audioFormat)

	if err != nil {
		return nil, err
	}

	return &Track{
		Name:      info.Title,
		Author:    info.Author,
		Duration:  info.Duration.Milliseconds(),
		ThumbURL:  sutils.Fmt("https://img.youtube.com/vi/%s/maxresdefault.jpg", info.ID),
		URL:       sutils.Fmt("https://youtu.be/%s", info.ID),
		StreamURL: resultUrl,
		isOpus:    isOpus,
	}, nil

}

func GetPlayer(guildId string) *Player {

	result, exist := players[guildId]

	if !exist {
		return nil
	}

	return result
}

func AddPlayer(player *Player) *Player {
	players[player.Guild.ID] = player
	return player
}
