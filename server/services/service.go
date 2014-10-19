package services

import (
	"apertoire.net/mediabase/server/model"
)

type Service interface {
	ConfigChanged(conf *model.Config)
}
