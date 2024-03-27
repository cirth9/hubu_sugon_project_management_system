package repo

import (
	"context"
	"test.com/project-project/internal/data"
)

type DepartmentRepo interface {
	FindDepartmentById(ctx context.Context, id int64) (*data.Department, error)
	FindDepartment(ctx context.Context, organizationCode int64, parentDepartmentCode int64, name string) (*data.Department, error)
	Save(dpm *data.Department) error
	ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, size int64) (list []*data.Department, total int64, err error)
}
