package message

import (
	"apertoire.net/mediabase/model"
)

type MovieScan struct {
	Payload *model.MovieScanReq
	Reply   chan *model.MovieScanRep
}

// type UserAuth struct {
// 	Payload *model.UserAuthReq
// 	Reply   chan *model.UserAuthRep
// }

// type UserData struct {
// 	Payload *model.UserDataReq
// 	Reply   chan *model.UserDataRep
// }
