package dao

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type FileDao struct {
	conn *gorms.GormConn
}

func (f *FileDao) FindByIds(ctx context.Context, ids []int64) (list []*data.File, err error) {
	session := f.conn.Session(ctx)
	err = session.Model(&data.File{}).Where("id in (?)", ids).Find(&list).Error
	return
}

func (f *FileDao) Save(ctx context.Context, file *data.File) error {
	err := f.conn.Session(ctx).Save(&file).Error
	return err
}

func NewFileDao() *FileDao {
	return &FileDao{
		conn: gorms.New(),
	}
}
