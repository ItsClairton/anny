package providers

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"

	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/buger/jsonparser"
)

type Song struct {
	Provider, Title, Uploader        string
	ThumbnailURL, DirectURL, PageURL string
	RawDuration                      string
	IsLive, IsOpus                   bool
}

func (s *Song) DisplayProvider() string {
	switch s.Provider {
	case "TwitchStream":
		return utils.Fmt("%s Twitch", emojis.Twitch)
	case "TwitchVod":
		return utils.Fmt("%s Twitch", emojis.Twitch)
	case "Youtube":
		return utils.Fmt("%s YouTube", emojis.Youtube)
	case "Soundcloud":
		return utils.Fmt("%s SoundCloud", emojis.Soundcloud)
	default:
		return s.Provider
	}
}

func (s *Song) Duration() string {
	if s.IsLive {
		return "--:--"
	}
	if s.RawDuration != "--:--" {
		return s.RawDuration
	}
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", s.DirectURL)

	buffer, err := cmd.Output()
	if err != nil {
		return "--:--"
	}

	duration, _ := jsonparser.GetString(buffer, "format", "duration")
	if floatDuration, err := strconv.ParseFloat(duration, 32); err == nil {
		s.RawDuration = utils.ToDisplayTime(floatDuration)
	}

	return s.RawDuration
}

func FindSong(argument string) (*Song, error) {
	cmd := exec.Command("youtube-dl", "--skip-download", "-f", "bestaudio/best", "--no-playlist", "--dump-json", argument)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	buffer, err := cmd.Output()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	provider, _ := jsonparser.GetString(buffer, "extractor_key")
	title, _ := jsonparser.GetString(buffer, "title")
	uploader, _ := jsonparser.GetString(buffer, "uploader")
	isLive, _ := jsonparser.GetBoolean(buffer, "is_live")
	thumbnailURL, _ := jsonparser.GetString(buffer, "thumbnail")
	directURL, _ := jsonparser.GetString(buffer, "url")
	pageURL, _ := jsonparser.GetString(buffer, "webpage_url")
	rawDuration, _ := jsonparser.GetInt(buffer, "duration")

	if provider == "TwitchStream" {
		title, _ = jsonparser.GetString(buffer, "description")
	}

	duration := "--:--"
	if rawDuration > 0 {
		duration = utils.ToDisplayTime(float64(rawDuration))
	} else {

	}

	return &Song{
		Provider:     provider,
		Title:        title,
		Uploader:     uploader,
		IsLive:       isLive,
		RawDuration:  duration,
		ThumbnailURL: thumbnailURL,
		DirectURL:    directURL,
		PageURL:      pageURL,
	}, nil
}
