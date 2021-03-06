package message

import (
	"apertoire.net/mediabase/server/model"
)

type ScanMovies struct {
	Reply chan string
}

type PruneMovies struct {
	Reply chan string
}

type MoviesDTO struct {
	Count  uint64   `json:"count"`
	Movies []*Movie `json:"movies"`
}

type Movies struct {
	Options Options
	Reply   chan *MoviesDTO
}

// type SearchMovies struct {
// 	Term  string
// 	Reply chan []*Movie
// }

type CheckMovie struct {
	Movie  *Movie
	Result chan bool
}

type SingleMovie struct {
	Movie *Movie
	Reply chan *Movie
}

type Movie struct {
	Id                   uint64  `json:"id"`
	Title                string  `json:"title"`
	Original_Title       string  `json:"original_title"`
	File_Title           string  `json:"file_title"`
	Year                 string  `json:"year"`
	Runtime              uint64  `json:"runtime"`
	Tmdb_Id              uint64  `json:"tmdb_id"`
	Imdb_Id              string  `json:"imdb_id"`
	Overview             string  `json:"overview"`
	Tagline              string  `json:"tagline"`
	Resolution           string  `json:"resolution"`
	FileType             string  `json:"filetype"`
	Location             string  `json:"location"`
	Cover                string  `json:"cover"`
	Backdrop             string  `json:"backdrop"`
	Genres               string  `json:"genres"`
	Vote_Average         float64 `json:"vote_average"`
	Vote_Count           uint64  `json:"vote_count"`
	Production_Countries string  `json:"production_countries"`
	Added                string  `json:"added"`
	Modified             string  `json:"modified"`
	Last_Watched         string  `json:"last_watched"`
	All_Watched          string  `json:"all_watched"`
	Count_Watched        uint64  `json:"count_watched"`
	Score                uint64  `json:"score"`
	Director             string  `json:"director"`
	Writer               string  `json:"writer"`
	Actors               string  `json:"actors"`
	Awards               string  `json:"awards"`
	Imdb_Rating          float64 `json:"imdb_rating"`
	Imdb_Votes           uint64  `json:"imdb_votes"`
}

type Media struct {
	BaseUrl       string
	SecureBaseUrl string
	BasePath      string
	Movie         *Movie
	Forced        bool
}

type Context struct {
	Message   string `json:"message"`
	Backdrop  string `json:"backdrop"`
	Completed bool   `json:"completed"`
}

type Status struct {
	Reply chan *Context
}

type GetConfig struct {
	Reply chan *model.Config
}

type SaveConfig struct {
	Config *model.Config
	Reply  chan *model.Config
}

type ConfigChanged struct {
	Config *model.Config
}

type Options struct {
	SearchTerm string `json:"searchTerm"`
	Current    uint64 `json:"current"`
	Limit      uint64 `json:"limit"`
	SortBy     string `json:"sortBy"`
	SortOrder  string `json:"sortOrder"`
	FilterBy   string `json:"filterBy"`
}
