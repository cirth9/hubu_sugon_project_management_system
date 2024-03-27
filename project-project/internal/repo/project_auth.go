package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type ProjectAuthRepo interface {
	FindAuthList(ctx context.Context, orgCode int64) (list []*data.ProjectAuth, err error)
	FindAuthListPage(ctx context.Context, orgCode int64, page int64, pageSize int64) (list []*data.ProjectAuth, total int64, err error)
}
