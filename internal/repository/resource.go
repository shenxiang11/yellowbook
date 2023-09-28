package repository

import (
	"context"
	"yellowbook/internal/domain"
	"yellowbook/internal/repository/dao"
)

type IResourceRepository interface {
	Create(ctx context.Context, domain domain.Resource, uploadUserId uint64) error
}

type ResourceRepository struct {
	dao dao.IResourceDao
}

func NewResourceRepository(dao dao.IResourceDao) IResourceRepository {
	return &ResourceRepository{dao: dao}
}

func (r *ResourceRepository) Create(ctx context.Context, resource domain.Resource, uploadUserId uint64) error {
	return r.dao.Insert(ctx, dao.Resource{
		Url:          resource.Url,
		Purpose:      resource.Purpose,
		Mimetype:     resource.Mimetype,
		UploadUserId: uploadUserId,
	})
}
