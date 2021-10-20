package audio

import (
	"errors"
	"regexp"
	"time"

	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/Pauloo27/searchtube"
	"github.com/kkdai/youtube/v2"
)

var (
	client        = &youtube.Client{}
	videoRegex    = regexp.MustCompile(`^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)
	hlsRegex      = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#,?&*//=]*)(.m3u8)\b([-a-zA-Z0-9@:%_\+.~#,?&//=]*))`)
	playlistRegex = regexp.MustCompile(`^https?:\/\/(www.youtube.com|youtube.com)\/playlist(.*)$`)
)

type YouTubeProvider struct{}

func (YouTubeProvider) Name() string {
	return utils.Fmt("%s %s", emojis.Youtube, "YouTube")
}

func (YouTubeProvider) IsValid(term string) bool {
	return videoRegex.MatchString(term) || !utils.URLRegex.MatchString(term) || playlistRegex.MatchString(term)
}

func (p YouTubeProvider) Find(term string) (*SongResult, error) {
	if playlistRegex.MatchString(term) {
		result, err := client.GetPlaylist(term)
		if err != nil {
			return nil, err
		}

		songResult := &SongResult{Songs: []*Song{}, IsFromSearch: false, IsFromPlaylist: true}
		playlist := &Playlist{
			Title:  result.Title,
			Author: result.Author,
			URL:    utils.Fmt("https://youtube.com/playlist?list=%s", result.ID),
		}

		for _, item := range result.Videos {
			songResult.Songs = append(songResult.Songs, &Song{
				Title:     item.Title,
				URL:       utils.Fmt("https://youtu.be/%s", item.ID),
				Thumbnail: utils.Fmt("https://img.youtube.com/vi/%s/mqdefault.jpg", item.ID),
				Author:    item.Author,
				Duration:  item.Duration,
				Playlist:  playlist,
				Provider:  &YouTubeProvider{},
			})

			playlist.Duration += item.Duration
		}

		return songResult, nil
	}

	if videoRegex.MatchString(term) {
		song, err := p.getSong(term, nil)
		if err != nil {
			return nil, err
		}

		return &SongResult{Songs: []*Song{song}, IsFromSearch: false, IsFromPlaylist: false}, nil
	}

	results, err := searchtube.Search(term, 5)
	if err != nil {
		return nil, err
	}

	if len(results) < 1 {
		return nil, nil
	}

	songResult := &SongResult{Songs: []*Song{}, IsFromSearch: true, IsFromPlaylist: false}

	for _, video := range results {
		var duration time.Duration
		if !video.Live {
			duration, _ = video.GetDuration()
		}

		songResult.Songs = append(songResult.Songs, &Song{
			Title:     video.Title,
			URL:       video.URL,
			Author:    video.Uploader,
			Thumbnail: video.Thumbnail,
			Duration:  duration,
			IsLive:    video.Live,
			Provider:  &YouTubeProvider{},
		})
	}

	return songResult, nil
}

func (p *YouTubeProvider) GetInfo(song *Song) (*Song, error) {
	return p.getSong(song.URL, song.Playlist)
}

func (YouTubeProvider) getSong(term string, playlist *Playlist) (*Song, error) {
	video, err := client.GetVideo(term)
	if err != nil {
		return nil, err
	}
	streamingURL := ""

	if video.HLSManifestURL != "" {
		if streamingURL, err = getLiveURL(video.HLSManifestURL); err != nil {
			return nil, err
		}
	} else {
		if streamingURL, err = client.GetStreamURL(video, video.Formats.FindByItag(140)); err != nil {
			return nil, err
		}
	}

	return &Song{
		Title:        video.Title,
		URL:          utils.Fmt("https://youtu.be/%s", video.ID),
		Author:       video.Author,
		Thumbnail:    video.Thumbnails[len(video.Thumbnails)-1].URL,
		StreamingURL: streamingURL,
		Duration:     video.Duration,
		IsLive:       video.HLSManifestURL != "",
		Provider:     &YouTubeProvider{},
		Playlist:     playlist,
	}, nil
}

func getLiveURL(url string) (string, error) {
	body, err := utils.GetFromWebString(url)
	if err != nil {
		return "", err
	}

	if hlsURL := hlsRegex.FindString(body); hlsURL != "" {
		return hlsURL, nil
	} else {
		return "", errors.New("no valid URL found within HLS")
	}
}
