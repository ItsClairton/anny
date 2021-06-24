package anilist

import (
	"sort"
	"strings"

	"github.com/ItsClairton/Anny/utils"
	"github.com/buger/jsonparser"
)

type MediaTitle struct {
	JP string `json:"romaji"`
	EN string `json:"english"`
}

type Media struct {
	Id        int         `json:"id"`
	IdMal     int         `json:"idMal"`
	Type      string      `json:"type"`
	Title     *MediaTitle `json:"title"`
	Synopsis  string      `json:"description"`
	SiteURL   string      `json:"siteUrl"`
	Status    string      `json:"status"`
	Format    string      `json:"format"`
	Season    string      `json:"season"`
	StartDate *utils.Date `json:"startDate"`
	Trailer   *Trailer    `json:"trailer"`
	EndDate   *utils.Date `json:"endDate"`
	Episodes  int         `json:"episodes"`
	Chapters  int         `json:"chapters"`
	Volumes   int         `json:"volumes"`
	Genres    []string    `json:"genres"`
	Cover     *CoverImage `json:"coverImage"`
	Banner    string      `json:"bannerImage"`
	Source    string      `json:"source"`
	Studio    struct {
		Node []*StudioNode `json:"nodes"`
	} `json:"studios"`
	Staff struct {
		Edge []*StaffEdge `json:"edges"`
		Node []*StaffNode `json:"nodes"`
	} `json:"staff"`
}

type Trailer struct {
	Id   string `json:"id"`
	Site string `json:"site"`
}

type CoverImage struct {
	ExtraLarge string `json:"extraLarge"` // Se não tiver disponível, automáticamente a API do AniList manda a imagem do tipo Large
	Color      string `json:"color"`
}

type StudioNode struct {
	Name              string `json:"name"`
	IsAnimationStudio bool   `json:"isAnimationStudio"`
}

type StaffEdge struct {
	Role string `json:"role"`
}

type StaffNode struct {
	Name struct {
		Full string `json:"full"`
	} `json:"name"`
}

type MALBasicInfo struct {
	Genres []string
	Score  float64
}

func (m *Media) GetArts() []string {

	var arts []string

	for i, entry := range m.Staff.Edge {
		if strings.Contains(strings.ToLower(entry.Role), "art") || strings.Contains(strings.ToLower(entry.Role), "illustration") {
			name := m.Staff.Node[i].Name.Full
			s := sort.SearchStrings(arts, name)
			if !(s < len(arts) && arts[s] == name) { // O AniList as vezes pode retornar duplicado
				arts = append(arts, name)
			}
		}
	}

	return arts
}

func (m *Media) GetCreator() string {

	var name string
	pattern := "original creator"

	if m.Type == "MANGA" {
		pattern = "story"
	}

	for i, entry := range m.Staff.Edge {
		if strings.Contains(strings.ToLower(entry.Role), pattern) {
			name = m.Staff.Node[i].Name.Full
		}
	}

	return name
}

func (m *Media) GetTrailerURL() string {
	if m.Trailer == nil {
		return ""
	}

	switch m.Trailer.Site {
	case "youtube":
		return utils.Fmt("https://www.youtube.com/watch?v=%s", m.Trailer.Id)
	case "dailymotion":
		return utils.Fmt("https://www.dailymotion.com/video/", m.Trailer.Id)
	default:
		return ""
	}

}

func (m *Media) GetDirectors() []string {

	var directors []string

	for i, entry := range m.Staff.Edge {
		if strings.EqualFold(entry.Role, "Director") {
			name := m.Staff.Node[i].Name.Full
			s := sort.SearchStrings(directors, name)

			if !(s < len(directors) && directors[s] == name) { // O AniList as vezes pode retornar duplicado
				directors = append(directors, name)
			}
		}
	}

	return directors
}

func (m *Media) GetAnimationStudios() []string {

	var studios []string

	for _, e := range m.Studio.Node {
		if e.IsAnimationStudio {
			studios = append(studios, e.Name)
		}
	}

	return studios
}

func (m *Media) GetType() int {

	switch m.Format {
	case "TV":
		return 0
	case "TV_SHORT":
		return 0
	case "MOVIE":
		return 1
	case "SPECIAL":
		return 2
	case "OVA":
		return 3
	case "ONA":
		return 4
	case "MUSIC":
		return 5
	default:
		return -1
	}

}

func (m *Media) GetSource() int {
	switch m.Source {
	case "ORIGINAL":
		return 0
	case "MANGA":
		return 1
	case "LIGHT_NOVEL":
		return 2
	case "VISUAL_NOVEL":
		return 3
	case "VIDEO_GAME":
		return 4
	case "OTHER":
		return 5
	case "NOVEL":
		return 6
	case "DOUJINSHI":
		return 7
	case "ANIME":
		return 8
	default:
		return -1
	}
}

func (m *Media) GetSeason() int {

	switch m.Season {
	case "WINTER":
		return 0
	case "SPRING":
		return 1
	case "SUMMER":
		return 2
	case "FALL":
		return 3
	default:
		return -1
	}

}

func (m *Media) GetStatus() int {
	switch m.Status {
	case "FINISHED":
		return 0
	case "RELEASING":
		return 1
	case "NOT_YET_RELEASED":
		return 2
	case "CANCELLED":
		return 3
	case "HIATUS":
		return 4
	default:
		return -1
	}
}

func (m *Media) GetBasicFromMAL() (*MALBasicInfo, error) {

	result, err := utils.GetFromWeb(utils.Fmt("https://api.jikan.moe/v3/%s/%d", strings.ToLower(m.Type), m.IdMal))

	if err != nil {
		return nil, err
	}

	var genres []string

	jsonparser.ArrayEach(result, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		if err != nil {
			return
		}

		name, err := jsonparser.GetString(value, "name")
		if err == nil {
			s := sort.SearchStrings(m.Genres, name)

			if !(s < len(m.Genres) && m.Genres[s] == name) { // O AniList as vezes pode retornar duplicado
				genres = append(genres, name)
			}
		}
	}, "genres")

	score, _ := jsonparser.GetFloat(result, "score")

	return &MALBasicInfo{
		Genres: genres,
		Score:  score,
	}, nil
}

func SearchMediaAsManga(title string) (*Media, error) {

	result, err := Get(Query{
		Query: `query ($search: String) {
			Media (search: $search, type: MANGA, isAdult: false) {
			  id, idMal, type, siteUrl
			  title { romaji english }
			  description(asHtml: true)
			  status(version: 2)
			  source(version: 2)
			  startDate { year, month, day }
			  endDate { year, month, day }
			  trailer { id site }
			  coverImage { extraLarge color }
			  bannerImage, chapters, volumes, genres
			  staff {
				nodes { name { full } }
				edges { role }
			  }
			}
		  }`,
		Variables: struct {
			Query string `json:"search"`
		}{Query: title},
	})

	if err != nil {
		return nil, err
	}

	return result.Media, nil
}

func GetMediaAsAnime(id int) (*Media, error) {

	variables := struct {
		ID int `json:"id"`
	}{
		ID: id,
	}

	result, err := Get(Query{
		Query:     "query ($id: Int) { Media (id: $id, type: ANIME) { id idMal type title { romaji english } description(asHtml: true), status, season, format, siteUrl, source startDate { year month day }, endDate { year month day }, trailer { id site } coverImage { extraLarge, color } bannerImage episodes genres studios { nodes { name isAnimationStudio } } staff { nodes { name { full } } edges { role } } } }",
		Variables: variables,
	})

	if err != nil {
		return nil, err
	}

	return result.Media, nil

}

func SearchMediaAsAnime(title string) (*Media, error) {

	variables := struct {
		Query string `json:"search"`
	}{
		Query: title,
	}

	result, err := Get(Query{
		Query:     "query ($search: String) { Media (search: $search, type: ANIME){ id idMal type title { romaji english } description(asHtml: true), status, season, format, siteUrl, source startDate { year month day }, endDate { year month day }, trailer { id site } coverImage { extraLarge, color } bannerImage episodes genres studios { nodes { name isAnimationStudio } } staff { nodes { name { full } } edges { role } } } }",
		Variables: variables,
	})

	if err != nil {
		return nil, err
	}

	return result.Media, nil

}
