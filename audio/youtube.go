package audio

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

var client = youtube.Client{}

func GetTrack(id string, req *discordgo.User) (*Track, error) {
	video, err := client.GetVideo(id)
	if err != nil {
		return nil, err
	}

	format := video.Formats.FindByItag(251)
	if format == nil {
		return nil, errors.New("opus audio format not found")
	}

	stream, err := client.GetStreamURL(video, format)
	if err != nil {
		return nil, err
	}

	return &Track{
		ID:           video.ID,
		Title:        video.Title,
		Author:       video.Author,
		Requester:    req,
		StreamingUrl: stream,
		IsOpus:       true,
	}, nil
}
