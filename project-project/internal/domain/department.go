package domain

import (
	"context"
	"test.com/project-common/errs"
	"test.com/project-project/internal/dao"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/repo"
	"test.com/project-project/pkg/model"
	"time"
)

type DepartmentDomain struct {
	departmentRepo repo.DepartmentRepo
}

func (d *DepartmentDomain) FindDepartmentById(id int64) (*data.Department, *errs.BError) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dp, err := d.departmentRepo.FindDepartmentById(c, id)
	if err != nil {
		return nil, model.DBError
	}
	return dp, nil
}

func (d *DepartmentDomain) List(organizationCode int64, parentDepartmentCode int64, page int64, size int64) ([]*data.DepartmentDisplay, int64, *errs.BError) {
	list, total, err := d.departmentRepo.ListDepartment(organizationCode, parentDepartmentCode, page, size)
	if err != nil {
		return nil, 0, model.DBError
	}
	var dList []*data.DepartmentDisplay
	for _, v := range list {
		dList = append(dList, v.ToDisplay())
	}
	return dList, total, nil
}

func (d *DepartmentDomain) Save(
	organizationCode int64,
	departmentCode int64,
	parentDepartmentCode int64,
	name string) (*data.DepartmentDisplay, *errs.BError) {

	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	dpm, err := d.departmentRepo.FindDepartment(c, organizationCode, parentDepartmentCode, name)
	if err != nil {
		return nil, model.DBError
	}
	if dpm == nil {
		dpm = &data.Department{
			Name:             name,
			OrganizationCode: organizationCode,
			CreateTime:       time.Now().UnixMilli(),
		}
		if parentDepartmentCode > 0 {
			dpm.Pcode = parentDepartmentCode
		}
		err := d.departmentRepo.Save(dpm)
		if err != nil {
			return nil, model.DBError
		}
		return dpm.ToDisplay(), nil
	}
	return dpm.ToDisplay(), nil
}

func NewDepartmentDomain() *DepartmentDomain {
	return &DepartmentDomain{
		departmentRepo: dao.NewDepartmentDao(),
	}
}
