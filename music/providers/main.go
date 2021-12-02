package music

import "time"

var availableProviders = []Provider{&YoutubeProvider{}}

type Provider interface {
	DisplayName() string
	IsSupported(term string) bool
	Find(term string) (*QueryResult, error)
	IsLoaded(*Song) bool
	Load(*Song) error
}

type QueryResult struct {
	Songs    []*Song
	Playlist *Playlist
}

type Playlist struct {
	Title, Author, URL string
	Duration           time.Duration
}

type Song struct {
	Title, Author, Thumbnail, URL, StreamingURL string

	Duration       time.Duration
	IsLive, IsOpus bool
	Expires        time.Time

	provider Provider
}

func (s *Song) IsLoaded() bool {
	return s.provider.IsLoaded(s)
}

func (s *Song) Load() error {
	if s.IsLoaded() {
		return nil
	}

	return s.provider.Load(s)
}

func (s *Song) Provider() string {
	return s.provider.DisplayName()
}

func FindSong(term string) (*QueryResult, error) {
	for _, provider := range availableProviders {
		if provider.IsSupported(term) {
			return provider.Find(term)
		}
	}

	return nil, nil
}
