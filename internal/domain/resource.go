package domain

import (
	"github.com/shenxiang11/yellowbook-proto/proto"
	"time"
)

type Resource struct {
	Id         uint64
	Url        string
	Purpose    proto.ResourcePurpose
	Mimetype   string
	CreateTime time.Time
	UpdateTime time.Time
	UploadUser *User
}
