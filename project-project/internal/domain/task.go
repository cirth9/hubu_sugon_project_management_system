package domain

import (
	"context"
	"test.com/project-common/errs"
	"test.com/project-common/kk"
	"test.com/project-project/config"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/repo"
	"test.com/project-project/pkg/model"
)

type TaskDomain struct {
	taskRepo repo.TaskRepo
}

func NewTaskDomain() *TaskDomain {
	return &TaskDomain{
		taskRepo: dao.NewTaskDao(),
	}
}

func (d *TaskDomain) FindProjectIdByTaskId(taskId int64) (int64, bool, *errs.BError) {
	config.SendLog(kk.Info("Find", "TaskDomain.FindProjectIdByTaskId", kk.FieldMap{
		"taskId": taskId,
	}))
	task, err := d.taskRepo.FindTaskById(context.Background(), taskId)
	if err != nil {
		config.SendLog(kk.Error(err, "TaskDomain.FindProjectIdByTaskId.taskRepo.FindTaskById", kk.FieldMap{
			"taskId": taskId,
		}))
		return 0, false, model.DBError
	}
	if task == nil {
		return 0, false, nil
	}
	return task.ProjectCode, true, nil
}
