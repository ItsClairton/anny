package anilist

import (
	"sort"
	"strings"

	"github.com/ItsClairton/Anny/utils/sutils"
)

type MediaTitle struct {
	JP string `json:"romaji"`
	EN string `json:"english"`
}

type Media struct {
	Id        int        `json:"id"`
	IdMal     int        `json:"idMal"`
	Title     MediaTitle `json:"title"`
	Synopsis  string     `json:"description"`
	SiteURL   string     `json:"siteUrl"`
	Status    string     `json:"status"`
	Format    string     `json:"format"`
	Season    string     `json:"season"`
	StartDate Date       `json:"startDate"`
	Trailer   Trailer    `json:"trailer"`
	EndDate   Date       `json:"endDate"`
	Episodes  int        `json:"episodes"`
	Genres    []string   `json:"genres"`
	Cover     CoverImage `json:"coverImage"`
	Banner    string     `json:"bannerImage"`
	IsAdult   bool       `json:"isAdult"`
	Source    string     `json:"source"`
	Studio    struct {
		Node []StudioNode `json:"nodes"`
	} `json:"studios"`
	Staff struct {
		Edge []StaffEdge `json:"edges"`
		Node []StaffNode `json:"nodes"`
	} `json:"staff"`
}

type Trailer struct {
	Id   string `json:"id"`
	Site string `json:"site"`
}

type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
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

func (m *Media) GetCreator() string {

	var name string

	for i, entry := range m.Staff.Edge {
		if strings.EqualFold(entry.Role, "Original Creator") {
			name = m.Staff.Node[i].Name.Full
		}
	}

	return name
}

func (m *Media) GetTrailerURL() string {

	switch m.Trailer.Site {
	case "youtube":
		return sutils.Fmt("https://www.youtube.com/watch?v=%s", m.Trailer.Id)
	case "dailymotion":
		return sutils.Fmt("https://www.dailymotion.com/video/", m.Trailer.Id)
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

func (m *Media) GetPrettyStartDate() string {
	return sutils.Fmt("%d %s de %d", m.StartDate.Day, sutils.ToPrettyMonth(m.StartDate.Month), m.StartDate.Year)
}

func (m *Media) GetPrettyEndDate() string {
	return sutils.Fmt("%d %s de %d", m.EndDate.Day, sutils.ToPrettyMonth(m.EndDate.Month), m.EndDate.Year)
}

func (m *Media) GetPrettySeason() string {

	switch m.Season {
	case "WINTER":
		return "Inverno"
	case "SPRING":
		return "Primavera"
	case "SUMMER":
		return "Verão"
	case "FALL":
		return "Outono"
	default:
		return "N/A"
	}

}

func (m *Media) GetPrettyFormat() string {

	switch m.Format {
	case "TV":
		return "TV"
	case "TV_SHORT":
		return "TV"
	case "MOVIE":
		return "Filme"
	case "SPECIAL":
		return "Especial"
	case "OVA":
		return "OVA"
	case "ONA":
		return "ONA"
	case "MUSIC":
		return "Música"
	default:
		return "N/A"
	}

}

func (m *Media) GetPrettyStatus() string {

	switch m.Status {
	case "FINISHED":
		return "FInalizado"
	case "RELEASING":
		return "Em Lançamento"
	case "NOT_YET_RELEASED":
		return "Não Lançado"
	case "CANCELLED":
		return "Cancelado"
	case "HIATUS":
		return "Pausado"
	default:
		return "N/A"
	}

}

func (m *Media) GetPrettySource() string {

	switch m.Source {
	case "ORIGINAL":
		return "Original"
	case "MANGA":
		return "Manga"
	case "LIGHT_NOVEL":
		return "Light Novel"
	case "VISUAL_NOVEL":
		return "Visual Novel"
	case "VIDEO_GAME":
		return "Jogos"
	case "OTHER":
		return "Outros"
	case "NOVEL":
		return "Novel"
	case "DOUJINSHI":
		return "Doujinshi"
	case "ANIME":
		return "Anime"
	default:
		return "N/A"
	}

}

func GetMediaAsAnime(id int) (Media, error) {

	variables := struct {
		ID int `json:"id"`
	}{
		ID: id,
	}

	result, err := Get(Query{
		Query:     "query ($id: Int) { Media (id: $id, type: ANIME) { id title { romaji english } description, status, season, format, siteUrl, source startDate { year month day }, endDate { year month day }, trailer { id site } coverImage { extraLarge, color } bannerImage isAdult episodes genres studios { nodes { name isAnimationStudio } } staff { nodes { name { full } } edges { role } } } }",
		Variables: variables,
	})

	if err != nil {
		return Media{}, err
	}

	return result.Media, nil

}

func SearchMediaAsAnime(title string) (Media, error) {

	variables := struct {
		Query string `json:"search"`
	}{
		Query: title,
	}

	result, err := Get(Query{
		Query:     "query ($search: String) { Media (search: $search, type: ANIME){ id title { romaji english } description, status, season, format, siteUrl, source startDate { year month day }, endDate { year month day }, trailer { id site } coverImage { extraLarge, color } bannerImage isAdult episodes genres studios { nodes { name isAnimationStudio } } staff { nodes { name { full } } edges { role } } } }",
		Variables: variables,
	})

	if err != nil {
		return Media{}, err
	}

	return result.Media, nil

}
