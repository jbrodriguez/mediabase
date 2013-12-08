package message

type ScanMovies struct {
	Reply chan string
}

type GetMovies struct {
	Reply chan []*Movie
}

type SearchMovies struct {
	Term  string
	Reply chan []*Movie
}

type Movie struct {
	Title          string `json: title`
	Original_Title string `json: `
	Year           string `json: year`
	Runtime        uint64 `json: runtime`
	Tmdb_Id        uint64 `json: tmdb_id`
	Imdb_Id        string `json: imdb_id`
	Overview       string `json: overview`
	Tagline        string `json: tagline`
	Resolution     string `json: resolution`
	FileType       string `json: filetype`
	Location       string `json: location`
	Cover          string `json: cover`
	Backdrop       string `json: backdrop`
}

type Media struct {
	BaseUrl       string `json: path`
	SecureBaseUrl string `json: id`
	Movie         *Movie `json: movie`
}

// type UserAuth struct {
// 	Payload *model.UserAuthReq
// 	Reply   chan *model.UserAuthRep
// }

// type UserData struct {
// 	Payload *model.UserDataReq
// 	Reply   chan *model.UserDataRep
// }
