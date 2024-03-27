package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type TaskWorkTimeRepo interface {
	Save(ctx context.Context, twt *data.TaskWorkTime) error
	FindWorkTimeList(ctx context.Context, taskCode int64) (list []*data.TaskWorkTime, err error)
}
