package audio

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kkdai/youtube/v2"
)

var client = youtube.Client{}

func GetTrack(id string, req *discordgo.User) (*Track, error) {
	video, err := client.GetVideo(id)
	if err != nil {
		return nil, err
	}

	stream, err := client.GetStreamURL(video, video.Formats.FindByItag(251))
	if err != nil {
		return nil, err
	}

	return &Track{
		ID:           video.ID,
		Name:         video.Title,
		Author:       video.Author,
		Requester:    req,
		StreamingUrl: stream,
		IsOpus:       true,
	}, nil
}
