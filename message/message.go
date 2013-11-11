package message

import (
	"apertoire.net/moviebase/model"
)

type UserAuth struct {
	Payload *model.UserAuthReq
	Reply   chan *model.UserAuthRep
}

type UserData struct {
	Payload *model.UserDataReq
	Reply   chan *model.UserDataRep
}
