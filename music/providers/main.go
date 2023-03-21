package music

import (
	"time"
)

var availableProviders = []Provider{&YoutubeProvider{}}

type Provider interface {
	DisplayName() string
	IsSupported(term string, query bool) bool
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
	Title     string `json:"title"`
	Author    string `json:"author"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
	MediaURL  string `json:"-"`

	Duration time.Duration `json:"duration"`
	IsLive   bool          `json:"isLive"`
	IsOpus   bool          `json:"-"`
	Expires  time.Time     `json:"-"`

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

func FindSong(term string, querySupport bool) (*QueryResult, error) {
	provider := FindByInput(term, querySupport)

	if provider == nil {
		return nil, nil
	}

	return provider.Find(term)
}

func FindByInput(term string, querySupport bool) Provider {
	for _, provider := range availableProviders {
		if provider.IsSupported(term, querySupport) {
			return provider
		}
	}

	return nil
}
