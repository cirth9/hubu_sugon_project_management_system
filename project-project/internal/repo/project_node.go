package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type ProjectNodeRepo interface {
	FindAll(ctx context.Context) (list []*data.ProjectNode, err error)
}
