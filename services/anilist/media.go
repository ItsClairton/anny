package anilist

type MediaTitle struct {
	JP string `json:"romaji"`
	EN string `json:"english"`
}

type Media struct {
	Id         int        `json:"id"`
	IdMal      int        `json:"idMal"`
	Title      MediaTitle `json:"title"`
	Synopsis   string     `json:"description"`
	Status     string     `json:"status"`
	Season     string     `json:"season"`
	SeasonYear int        `json:"seasonYear"`
	StartDate  Date       `json:"startDate"`
	Trailer    Trailer    `json:"trailer"`
	EndDate    Date       `json:"endDate"`
	Episodes   int        `json:"episodes"`
	Genres     []string   `json:"genres"`
	Cover      CoverImage `json:"coverImage"`
	Banner     string     `json:"bannerImage"`
	IsAdult    bool       `json:"isAdult"`
	Source     string     `json:"source"`
	Studio     struct {
		Node StudioNode `json:"node"`
	} `json:"studios"`
	Staff struct {
		Edge StaffEdge `json:"edge"`
		Node StaffNode `json:"node"`
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

func GetMediaAsAnime(id int) (Media, error) {

	variables := struct {
		ID int `json:"id"`
	}{
		ID: id,
	}

	result, err := Get(Query{
		Query:     "query ($id: Int) { Media (id: $id, type: ANIME) { id title { romaji english } description, status, season, seasonYear, source startDate { year month day }, endDate { year month day }, trailer { id site } coverImage { extraLarge } bannerImage isAdult episodes genres studios { nodes { name isAnimationStudio } } staff { nodes { name { full } } edges { role } } } }",
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
		Query:     "query ($search: String) { Media (search: $search, type: ANIME) { id title { romaji english } description, status, season, seasonYear, source startDate { year month day }, endDate { year month day }, trailer { id site } coverImage { extraLarge } bannerImage isAdult episodes genres studios { nodes { name isAnimationStudio } } staff { nodes { name { full } } edges { role } } } }",
		Variables: variables,
	})

	if err != nil {
		return Media{}, err
	}

	return result.Media, nil

}
