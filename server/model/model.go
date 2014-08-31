package model

type ScanMovieReq struct {
	Start bool `json: start`
}

type ScanMovieRep struct {
	Started bool `json: started`
}

// type UserAuthReq struct {
// 	Email    string `json: email`
// 	Password string `json: password`
// }

// type UserAuthRep struct {
// 	Id    int8   `json: id`
// 	Email string `json: email`
// }

// type UserDataReq struct {
// 	Id int8 `json: id`
// }

// type UserDataRep struct {
// 	Name  string `json: name`
// 	Email string `json: email`
// }
