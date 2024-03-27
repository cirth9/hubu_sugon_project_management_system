package repo

import (
	"context"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/database"
)

type MemberRepo interface {
	GetMemberByEmail(ctx context.Context, email string) (bool, error)
	GetMemberByAccount(ctx context.Context, account string) (bool, error)
	GetMemberByMobile(ctx context.Context, mobile string) (bool, error)
	SaveMember(conn database.DbConn, ctx context.Context, mem *member.Member) error
	FindMember(ctx context.Context, account string, pwd string) (mem *member.Member, err error)
	FindMemberById(background context.Context, id int64) (mem *member.Member, err error)
	FindMemberByIds(background context.Context, ids []int64) (list []*member.Member, err error)
	UpdateMember(background context.Context, mem1 *member.Member) (mem *member.Member, err error)
}
