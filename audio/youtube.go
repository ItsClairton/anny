package audio

import (
	"errors"
	"strings"
	"time"

	"github.com/ItsClairton/Anny/utils"
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
	isOpus := false

	if format == nil {
		format = &video.Formats.WithAudioChannels()[0]

		if format == nil {
			return nil, errors.New("opus audio format not found")
		}
		isOpus = strings.Contains(format.MimeType, "opus")
	}

	stream, err := client.GetStreamURL(video, format)
	if err != nil {
		return nil, err
	}

	return &Track{
		Title:        video.Title,
		ID:           video.ID,
		URL:          utils.Fmt("https://youtu.be/%s", video.ID),
		ThumbnailUrl: video.Thumbnails[len(video.Thumbnails)-1].URL,
		Author:       video.Author,
		Requester:    req,
		Duration:     video.Duration,
		IsOpus:       isOpus,
		StreamingUrl: stream,
	}, nil
}

func GetPlaylist(id string, req *discordgo.User) ([]*Track, time.Duration, error) {
	playlist, err := client.GetPlaylist(id)
	if err != nil {
		return nil, 0, err
	}

	tracks := []*Track{}
	var duration time.Duration

	for _, video := range playlist.Videos {
		duration = duration + video.Duration
		tracks = append(tracks, &Track{
			Title:        video.Title,
			ID:           video.ID,
			URL:          utils.Fmt("https://youtu.be/%s", video.ID),
			ThumbnailUrl: utils.Fmt("https://img.youtube.com/vi/%s/mqdefault.jpg", video.ID),
			Author:       video.Author,
			Duration:     video.Duration,
			Requester:    req,
			Playlist: &Playlist{
				ID:     playlist.ID,
				Title:  playlist.Title,
				Author: playlist.Author,
			},
		})
	}

	return tracks, duration, nil
}
