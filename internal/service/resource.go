package service

import (
	"context"
	"github.com/shenxiang11/yellowbook-proto/proto"
	"github.com/spf13/viper"
	"mime/multipart"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository"
	"yellowbook/internal/service/oss"
	"yellowbook/pkg/logger"
)

type IResourceService interface {
	Upload(ctx context.Context, f *multipart.FileHeader, purpose proto.ResourcePurpose, uid uint64) (string, error)
	GetResourceCategoryList() any
}

type ResourceService struct {
	ossSrv oss.IService
	repo   repository.IResourceRepository
	l      logger.Logger
}

func NewResourceService(ossSrv oss.IService, repo repository.IResourceRepository, l logger.Logger) IResourceService {
	return &ResourceService{
		ossSrv: ossSrv,
		repo:   repo,
		l:      l,
	}
}

func (s *ResourceService) Upload(ctx context.Context, f *multipart.FileHeader, purpose proto.ResourcePurpose, uid uint64) (string, error) {
	url, err := s.ossSrv.Upload(f)
	if err != nil {
		return "", err
	}

	mimeType := f.Header.Get("Content-Type")

	err = s.repo.Create(ctx, domain.Resource{
		Url:      url,
		Purpose:  purpose,
		Mimetype: mimeType,
	}, uid)
	if err != nil {
		s.l.Warn("OSS 上传成功，系统记录失败", logger.Field{
			Key:   "url",
			Value: url,
		})
	}

	return url, nil
}

func (s *ResourceService) GetResourceCategoryList() any {
	c := viper.Get("dict_resource_type")
	return c
}
