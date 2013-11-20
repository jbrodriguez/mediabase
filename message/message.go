package message

type ScanMovies struct {
	Reply chan string
}

type GetMovies struct {
	Reply chan []Movie
}

type Movie struct {
	Resolution string `json: resolution`
	Name       string `json: name`
	Year       string `json: year`
	Type       string `json: type`
	Path       string `json: path`
	Picture    string `json: picture`
}

type Picture struct {
	Path string `json: path`
	Id   string `json: id`
}

// type UserAuth struct {
// 	Payload *model.UserAuthReq
// 	Reply   chan *model.UserAuthRep
// }

// type UserData struct {
// 	Payload *model.UserDataReq
// 	Reply   chan *model.UserDataRep
// }
