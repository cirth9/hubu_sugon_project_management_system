package dao

import (
	"test.com/project-user/internal/database"
	"test.com/project-user/internal/database/gorms"
)

type TransactionImpl struct {
	conn database.DbConn
}

func (t *TransactionImpl) Action(f func(conn database.DbConn) error) error {
	t.conn.Begin()
	err := f(t.conn)
	if err != nil {
		t.conn.Rollback()
		return err
	}
	t.conn.Commit()
	return nil
}

func NewTransaction() *TransactionImpl {
	return &TransactionImpl{
		conn: gorms.NewTran(),
	}
}
