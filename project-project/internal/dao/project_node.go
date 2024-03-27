package dao

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database/gorms"
)

type ProjectNodeDao struct {
	conn *gorms.GormConn
}

func (m *ProjectNodeDao) FindAll(ctx context.Context) (pms []*data.ProjectNode, err error) {
	session := m.conn.Session(ctx)
	err = session.Model(&data.ProjectNode{}).Find(&pms).Error
	return
}

func NewProjectNodeDao() *ProjectNodeDao {
	return &ProjectNodeDao{
		conn: gorms.New(),
	}
}
