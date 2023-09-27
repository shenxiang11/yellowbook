package ioc

import (
	"yellowbook/internal/service/oss"
)

func InitOss() oss.IService {
	return oss.NewService()
}
