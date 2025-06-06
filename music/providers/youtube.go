package music

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ItsClairton/Anny/utils"
	"github.com/ItsClairton/Anny/utils/emojis"
	"github.com/Pauloo27/searchtube"
	"github.com/kkdai/youtube/v2"
	"github.com/pkg/errors"
)

type YoutubeProvider struct{}

var (
	videoRegex    = regexp.MustCompile(`^((?:https?:)?\/\/)?((?:www|m)\.)?((?:youtube\.com|youtu.be))(\/(?:[\w\-]+\?v=|embed\/|v\/)?)([\w\-]+)(\S+)?$`)
	playlistRegex = regexp.MustCompile(`[&?]list=([A-Za-z0-9_-]{13,42})(&.*)?$`)
	hlsRegex      = regexp.MustCompile(`(https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#,?&*//=]*)(.m3u8)\b([-a-zA-Z0-9@:%_\+.~#,?&//=]*))`)

	client = &youtube.Client{}
	cache  = make(map[string]*Song)
)

func (YoutubeProvider) DisplayName() string {
	return utils.Fmt("%s YouTube", emojis.Youtube)
}

func (YoutubeProvider) IsSupported(term string, query bool) bool {
	return videoRegex.MatchString(term) || (query && !utils.LinkRegex.MatchString(term)) || playlistRegex.MatchString(term)
}

func (YoutubeProvider) IsLoaded(s *Song) bool {
	return !time.Now().Add(s.Duration).After(s.Expires)
}

func (provider YoutubeProvider) Load(s *Song) error {
	loadedSong, err := provider.handleVideo(s.URL)
	if err != nil {
		return err
	}

	s.MediaURL, s.IsOpus, s.Expires = loadedSong.MediaURL, loadedSong.IsOpus, loadedSong.Expires
	s.Thumbnail = loadedSong.Thumbnail
	return nil
}

func (provider YoutubeProvider) Find(term string) (*QueryResult, error) {
	if playlistRegex.MatchString(term) {
		return provider.handlePlaylist(term)
	}

	if videoRegex.MatchString(term) {
		if video, err := provider.handleVideo(term); err != nil {
			return nil, err
		} else {
			return &QueryResult{Songs: []*Song{video}}, nil
		}
	}

	items, err := searchtube.Search(term, 5)
	if err != nil {
		return nil, err
	}

	if len(items) < 1 {
		return nil, nil
	}

	result := &QueryResult{}
	for _, video := range items {
		duration, _ := video.GetDuration()

		result.Songs = append(result.Songs, &Song{
			Title:     video.Title,
			URL:       video.URL,
			Author:    video.Uploader,
			Thumbnail: video.Thumbnail,
			Duration:  duration,
			IsLive:    video.Live,
			provider:  provider,
		})
	}

	return result, nil
}

func (provider YoutubeProvider) handlePlaylist(URL string) (*QueryResult, error) {
	playlist, err := client.GetPlaylist(URL)
	if err != nil {
		if video, err := provider.handleVideo(URL); err == nil {
			return &QueryResult{Songs: []*Song{video}}, nil
		}

		return nil, err
	}

	result := &QueryResult{
		Songs: make([]*Song, len(playlist.Videos)),
		Playlist: &Playlist{
			Title:  playlist.Title,
			Author: playlist.Author,
			URL:    utils.Fmt("https://youtube.com/playlist?list=%s", playlist.ID),
		},
	}

	for i, item := range playlist.Videos {
		result.Playlist.Duration += item.Duration

		result.Songs[i] = &Song{
			Title:     item.Title,
			Author:    item.Author,
			Duration:  item.Duration,
			Thumbnail: utils.Fmt("https://img.youtube.com/vi/%s/mqdefault.jpg", item.ID),
			URL:       utils.Fmt("https://youtu.be/%s", item.ID),
			provider:  provider,
		}
	}

	return result, nil
}

func (provider YoutubeProvider) handleVideo(term string) (song *Song, err error) {
	if term, err = youtube.ExtractVideoID(term); err != nil {
		return nil, err
	}

	if cached := cache[term]; cached != nil && provider.IsLoaded(cached) {
		return cached, nil
	}

	video, err := client.GetVideo(term)
	if err != nil {
		return nil, err
	}

	mediaURL, isOpus := "", false
	if video.Duration > 0 {
		var formatList youtube.FormatList

		if formatList, isOpus = video.Formats.Itag(251), true; len(formatList) == 0 { // Opus
			formatList, isOpus = video.Formats.Itag(140), false // M4a
		}

		if mediaURL, err = client.GetStreamURL(video, &formatList[0]); err != nil {
			return nil, err
		}
	} else {
		if mediaURL, err = getLiveURL(video.HLSManifestURL); err != nil {
			return nil, err
		}
	}

	expires := time.Now().Add(10 * time.Minute)
	if video.Duration > 0 {
		if expires, err = getExpires(mediaURL); err != nil {
			return nil, err
		}
	}

	thumbnail := utils.Fmt("https://img.youtube.com/vi/%s/mqdefault.jpg", video.ID)
	if format := video.Formats.Quality("720p"); format != nil {
		thumbnail = utils.Fmt("https://img.youtube.com/vi/%s/maxresdefault.jpg", video.ID)
	}

	song = &Song{
		Title:     video.Title,
		Author:    video.Author,
		URL:       utils.Fmt("https://youtu.be/%s", video.ID),
		Duration:  video.Duration,
		Thumbnail: thumbnail,
		MediaURL:  mediaURL,
		Expires:   expires,
		IsLive:    video.Duration == 0,
		IsOpus:    isOpus,
		provider:  &provider,
	}

	if !song.IsLive {
		cache[term] = song
	}

	return song, nil
}

func getExpires(URL string) (time.Time, error) {
	firstIndex := strings.Index(URL, "expire=")
	if firstIndex == -1 {
		return time.Time{}, errors.New("unexpected URL")
	}

	endIndex := strings.Index(URL[firstIndex:], "&")
	if endIndex == -1 {
		return time.Time{}, errors.New("unexpected URL")
	}

	expires, err := strconv.Atoi(URL[firstIndex+7 : firstIndex+endIndex])
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(expires), 0), nil
}

func getLiveURL(manifestURL string) (string, error) {
	body, err := utils.FromWebString(manifestURL)
	if err != nil {
		return "", err
	}

	if hlsURL := hlsRegex.FindString(body); hlsURL != "" {
		return hlsURL, nil
	} else {
		return "", errors.New("no valid URL found within HLS")
	}
}
