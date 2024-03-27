package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type SourceLinkRepo interface {
	Save(ctx context.Context, link *data.SourceLink) error
	FindByTaskCode(ctx context.Context, taskCode int64) (list []*data.SourceLink, err error)
}
