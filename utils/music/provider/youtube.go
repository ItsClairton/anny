package provider

import (
	"errors"
	"net/url"

	"github.com/ItsClairton/Anny/utils"
	"github.com/buger/jsonparser"
	"github.com/kkdai/youtube/v2"
)

type YouTubeProvider struct{}

var client = youtube.Client{}

func (p YouTubeProvider) GetInfo(content string) (*PartialInfo, error) {

	result, err := utils.GetFromWeb(utils.Fmt("https://youtube-scrape.herokuapp.com/api/search?q=%s", url.QueryEscape(content)))
	if err != nil {
		return nil, err
	}

	id, err := jsonparser.GetUnsafeString(result, "results", "[0]", "video", "id")
	if err != nil {
		return nil, nil
	}

	title, _ := jsonparser.GetUnsafeString(result, "results", "[0]", "video", "title")
	author, _ := jsonparser.GetUnsafeString(result, "results", "[0]", "uploader", "username")
	duration, _ := jsonparser.GetUnsafeString(result, "results", "[0]", "video", "duration")

	return &PartialInfo{
		ID:       id,
		Title:    title,
		Author:   author,
		Duration: duration,
		URL:      utils.Fmt("https://youtu.be/%s", id),
		ThumbURL: utils.Fmt("https://img.youtube.com/vi/%s/maxresdefault.jpg", id),
		Provider: p,
	}, nil

}

func (p YouTubeProvider) GetStream(info *PartialInfo) (*StreamInfo, error) {

	video, err := client.GetVideo(info.ID)
	if err != nil {
		return nil, err
	}

	var format youtube.Format
	var isOpus bool
	if len(video.Formats.Type("opus")) < 1 {
		if len(video.Formats.Type("audio")) < 1 {
			return nil, errors.New("no avaliable audio formats")
		}

		format = video.Formats.Type("audio")[0]
		isOpus = false
	} else {
		format = video.Formats.Type("opus")[0]
		isOpus = true
	}

	streamUrl, err := client.GetStreamURL(video, &format)
	if err != nil {
		return nil, err
	}

	return &StreamInfo{
		IsOpus:    isOpus,
		StreamURL: streamUrl,
	}, nil
}
