package dao

import (
	"context"
	"test.com/project-project/internal/data"
	"test.com/project-project/internal/database"
	"test.com/project-project/internal/database/gorms"
)

type ProjectAuthNodeDao struct {
	conn *gorms.GormConn
}

func (p *ProjectAuthNodeDao) FindNodeStringList(ctx context.Context, authId int64) (list []string, err error) {
	session := p.conn.Session(ctx)
	err = session.Model(&data.ProjectAuthNode{}).Where("auth=?", authId).Select("node").Find(&list).Error
	return
}

func (p *ProjectAuthNodeDao) DeleteByAuthId(ctx context.Context, conn database.DbConn, authId int64) error {
	p.conn = conn.(*gorms.GormConn)
	tx := p.conn.Tx(ctx)
	err := tx.Where("auth=?", authId).Delete(&data.ProjectAuthNode{}).Error
	return err
}

func (p *ProjectAuthNodeDao) Save(ctx context.Context, conn database.DbConn, authId int64, nodes []string) error {
	p.conn = conn.(*gorms.GormConn)
	tx := p.conn.Tx(ctx)
	var list []*data.ProjectAuthNode
	for _, v := range nodes {
		pn := &data.ProjectAuthNode{
			Auth: authId,
			Node: v,
		}
		list = append(list, pn)
	}
	err := tx.Create(list).Error
	return err
}

func NewProjectAuthNodeDao() *ProjectAuthNodeDao {
	return &ProjectAuthNodeDao{
		conn: gorms.New(),
	}
}
