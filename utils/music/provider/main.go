package provider

type PartialInfo struct {
	Title, Duration, Author string
	ID, URL, ThumbURL       string
	Provider                Provider
}

type StreamInfo struct {
	StreamURL string
	IsOpus    bool
}

type Provider interface {
	GetInfo(string) (*PartialInfo, error)
	GetStream(*PartialInfo) (*StreamInfo, error)
}

var (
	YouProvider = YouTubeProvider{}
)

func GetInfo(content string) (*PartialInfo, error) {
	return YouProvider.GetInfo(content)
}
