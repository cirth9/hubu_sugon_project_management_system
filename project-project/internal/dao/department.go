package dao

import (
	"context"
	"gorm.io/gorm"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type DepartmentDao struct {
	conn *gorms.GormConn
}

func (d *DepartmentDao) Save(dpm *data.Department) error {
	err := d.conn.Session(context.Background()).Save(&dpm).Error
	return err
}

func (d *DepartmentDao) FindDepartment(ctx context.Context, organizationCode int64, parentDepartmentCode int64, name string) (*data.Department, error) {
	session := d.conn.Session(ctx)
	session = session.Model(&data.Department{}).Where("organization_code=? AND name=?", organizationCode, name)
	if parentDepartmentCode > 0 {
		session = session.Where("pcode=?", parentDepartmentCode)
	}
	var dp *data.Department
	err := session.Limit(1).Take(&dp).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return dp, err
}

func (d *DepartmentDao) ListDepartment(organizationCode int64, parentDepartmentCode int64, page int64, size int64) (list []*data.Department, total int64, err error) {
	session := d.conn.Session(context.Background())
	session = session.Model(&data.Department{})
	session = session.Where("organization_code=?", organizationCode)
	if parentDepartmentCode > 0 {
		session = session.Where("pcode=?", parentDepartmentCode)
	}
	err = session.Count(&total).Error
	err = session.Limit(int(size)).Offset(int((page - 1) * size)).Find(&list).Error
	return
}

func (d *DepartmentDao) FindDepartmentById(ctx context.Context, id int64) (dt *data.Department, err error) {
	session := d.conn.Session(ctx)
	err = session.Where("id=?", id).Find(&dt).Error
	return
}

func NewDepartmentDao() *DepartmentDao {
	return &DepartmentDao{
		conn: gorms.New(),
	}
}
