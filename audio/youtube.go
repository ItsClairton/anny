package audio

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
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
	playlistRegex = regexp.MustCompile(`[&?]list=([A-Za-z0-9_-]{18,42})(&.*)?$`)
	cache         = make(map[string]*Song)
)

type YouTubeProvider struct{}

func (*YouTubeProvider) Name() string {
	return utils.Fmt("%s %s", emojis.Youtube, "YouTube")
}

func (*YouTubeProvider) IsCompatible(term string) bool {
	return videoRegex.MatchString(term) || !utils.URLRegex.MatchString(term) || playlistRegex.MatchString(term)
}

func (p *YouTubeProvider) Find(term string) (*SongResult, error) {
	if playlistRegex.MatchString(term) {
		return p.getPlaylist(term)
	}

	if videoRegex.MatchString(term) {
		song, err := p.getSong(term, nil, 1)
		if err != nil {
			return nil, err
		}

		return &SongResult{Songs: []*Song{song}, IsFromSearch: false, IsFromPlaylist: false}, nil
	}

	items, err := searchtube.Search(term, 1)
	if err != nil {
		return nil, err
	}

	if len(items) < 1 {
		return nil, nil
	}

	result := &SongResult{Songs: []*Song{}, IsFromSearch: true, IsFromPlaylist: false}

	for _, video := range items {
		var duration time.Duration
		if !video.Live {
			duration, _ = video.GetDuration()
		}

		result.Songs = append(result.Songs, &Song{
			Title:     video.Title,
			URL:       video.URL,
			Author:    video.Uploader,
			Thumbnail: video.Thumbnail,
			Duration:  duration,
			IsLive:    video.Live,
			provider:  p,
		})
	}

	return result, nil
}

func (p *YouTubeProvider) Load(song *Song) (*Song, error) {
	return p.getSong(song.URL, song.Playlist, 1)
}

func (p *YouTubeProvider) getPlaylist(term string) (*SongResult, error) {
	data, err := client.GetPlaylist(term)
	if err != nil {
		return nil, err
	}

	result := &SongResult{Songs: []*Song{}, IsFromSearch: false, IsFromPlaylist: true}

	playlist := &Playlist{
		Title: data.Title, Author: data.Author,
		URL: utils.Fmt("https://youtube.com/playlist?list=%s", data.ID),
	}

	for _, item := range data.Videos {
		playlist.Duration += item.Duration

		result.Songs = append(result.Songs, &Song{
			Title:     item.Title,
			URL:       utils.Fmt("https://youtu.be/%s", item.ID),
			Thumbnail: utils.Fmt("https://img.youtube.com/vi/%s/mqdefault.jpg", item.ID),
			Duration:  item.Duration,
			Playlist:  playlist,
			provider:  p,
		})
	}

	return result, nil
}

func (p *YouTubeProvider) getSong(term string, playlist *Playlist, attempts int) (*Song, error) {
	cached := cache[term]
	if cached != nil && p.IsLoaded(cached) {
		return cached, nil
	}

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
		format := video.Formats.FindByItag(251) // Opus
		if format == nil {
			format = video.Formats.FindByItag(140) // M4a
		}

		if streamingURL, err = client.GetStreamURL(video, format); err != nil {
			return nil, err
		}
	}

	res, err := http.Get(streamingURL)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	if res.StatusCode >= 400 {
		if attempts >= 5 {
			return nil, fmt.Errorf("the server responded with unexpected %d status code after 5 attempts", res.StatusCode)
		}

		attempts++
		return p.getSong(term, playlist, attempts)
	}

	song := &Song{
		Title:        video.Title,
		URL:          utils.Fmt("https://youtu.be/%s", video.ID),
		Author:       video.Author,
		Thumbnail:    video.Thumbnails[len(video.Thumbnails)-1].URL,
		StreamingURL: streamingURL,
		Duration:     video.Duration,
		IsLive:       video.HLSManifestURL != "",
		Playlist:     playlist,
		provider:     p,
	}

	expiresStr := res.Request.URL.Query().Get("expire")
	if expiresStr != "" {
		if expires, err := strconv.Atoi(expiresStr); err == nil {
			song.Expires = time.Unix(int64(expires), 0)
		}
	}

	if !song.IsLive {
		cache[song.URL] = song
	}
	return song, nil
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

func (YouTubeProvider) IsLoaded(song *Song) bool {
	return !song.Expires.IsZero() && !time.Now().Add(song.Duration).After(song.Expires)
}
