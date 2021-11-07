package audio

import "time"

var availableProviders = []SongProvider{&YouTubeProvider{}}

type Song struct {
	Title, Author, Thumbnail, URL, StreamingURL string

	Duration time.Duration
	Playlist *Playlist
	IsLive   bool

	provider SongProvider
}

type SongProvider interface {
	Name() string
	IsValid(string) bool
	Find(string) (*SongResult, error)
	Load(*Song) (*Song, error)
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

func (s *Song) IsLoaded() bool {
	return s.StreamingURL != ""
}

func (s *Song) Load() (*Song, error) {
	return s.provider.Load(s)
}

func (s *Song) Provider() string {
	return s.provider.Name()
}
