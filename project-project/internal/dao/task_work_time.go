package dao

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type TaskWorkTimeDao struct {
	conn *gorms.GormConn
}

func (t *TaskWorkTimeDao) Save(ctx context.Context, twt *data.TaskWorkTime) error {
	session := t.conn.Session(ctx)
	err := session.Save(&twt).Error
	return err
}

func (t *TaskWorkTimeDao) FindWorkTimeList(ctx context.Context, taskCode int64) (list []*data.TaskWorkTime, err error) {
	session := t.conn.Session(ctx)
	err = session.Model(&data.TaskWorkTime{}).Where("task_code=?", taskCode).Find(&list).Error
	return
}

func NewTaskWorkTimeDao() *TaskWorkTimeDao {
	return &TaskWorkTimeDao{
		conn: gorms.New(),
	}
}
