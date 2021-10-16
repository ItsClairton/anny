package audio

import "time"

var availableProviders = []SongProvider{&YouTubeProvider{}}

type Song struct {
	Title, Author, Thumbnail, URL, StreamingURL string

	Duration time.Duration
	IsLive   bool

	Provider SongProvider
	Playlist *Playlist
}

type SongProvider interface {
	PrettyName() string
	IsValid(string) bool
	Find(string) (*SongResult, error)
	GetInfo(*Song) (*Song, error)
}

type SongResult struct {
	Songs []*Song

	IsFromSearch   bool
	IsFromPlaylist bool
}

type Playlist struct {
	Title, Author, URL string
	Duration           time.Duration
}

func FindSong(term string) (*SongResult, error) {
	for _, provider := range availableProviders {
		if provider.IsValid(term) {
			return provider.Find(term)
		}
	}
	return nil, nil
}
