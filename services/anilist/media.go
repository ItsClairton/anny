package anilist

import (
	"sort"
	"strings"

	"github.com/ItsClairton/Anny/utils/rest"
	"github.com/ItsClairton/Anny/utils/sutils"
	"github.com/buger/jsonparser"
)

type MediaTitle struct {
	JP string `json:"romaji"`
	EN string `json:"english"`
}

type Media struct {
	Id        int        `json:"id"`
	IdMal     int        `json:"idMal"`
	Type      string     `json:"type"`
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
	Chapters  int        `json:""`
	Volumes   int        `json:""`
	Genres    []string   `json:"genres"`
	Cover     CoverImage `json:"coverImage"`
	Banner    string     `json:"bannerImage"`
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

func (m *Media) GetArts() []string {

	var arts []string

	for i, entry := range m.Staff.Edge {
		if strings.Contains(strings.ToLower(entry.Role), "art") {
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

	if m.Status == "NOT_YET_RELEASED" {
		if m.StartDate.Year > 0 {
			return sutils.Fmt("Previsto para %d", m.StartDate.Year)
		} else {
			return "Ainda não divulgado."
		}
	}

	if m.StartDate.Day == 0 || m.StartDate.Month == 0 || m.StartDate.Year == 0 {
		return "N/A"
	}

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

func (m *Media) GetScoreFromMAL() (float64, error) {
	result, err := rest.Get(sutils.Fmt("https://api.jikan.moe/v3/%s/%d", strings.ToLower(m.Type), m.IdMal))

	if err != nil {
		return -1, nil
	}

	return jsonparser.GetFloat(result, "score")
}

func SearchMediaAsManga(title string) (Media, error) {

	result, err := Get(Query{
		Query: `query ($search: String) {
			Media (search: $search, type: MANGA) {
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
		return Media{}, err
	}

	return result.Media, nil
}

func GetMediaAsAnime(id int) (Media, error) {

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
		Query:     "query ($search: String) { Media (search: $search, type: ANIME){ id idMal type title { romaji english } description(asHtml: true), status, season, format, siteUrl, source startDate { year month day }, endDate { year month day }, trailer { id site } coverImage { extraLarge, color } bannerImage episodes genres studios { nodes { name isAnimationStudio } } staff { nodes { name { full } } edges { role } } } }",
		Variables: variables,
	})

	if err != nil {
		return Media{}, err
	}

	return result.Media, nil

}
