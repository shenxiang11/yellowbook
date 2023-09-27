package service

import (
	"github.com/spf13/viper"
	"mime/multipart"
	"yellowbook/internal/service/oss"
)

type IResourceService interface {
	Upload(f *multipart.FileHeader, uid uint64) (string, error)
	GetResourceCategoryList() any
}

type ResourceService struct {
	ossSrv oss.IService
}

func NewResourceService(ossSrv oss.IService) IResourceService {
	return &ResourceService{
		ossSrv: ossSrv,
	}
}

// C 端上传
func (s *ResourceService) Upload(f *multipart.FileHeader, uid uint64) (string, error) {
	url, err := s.ossSrv.Upload(f)
	if err != nil {
		return "", err
	}

	//mimeType := f.Header.Get("Content-Type")

	return url, nil
}

func (s *ResourceService) GetResourceCategoryList() any {
	c := viper.Get("dict_resource_type")
	return c
}
