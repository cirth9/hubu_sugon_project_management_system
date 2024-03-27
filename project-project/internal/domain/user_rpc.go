package domain

import (
	"context"
	"test.com/project-grpc/user/login"
	"test.com/project-project/internal/rpc"
	"time"
)

type UserRpcDomain struct {
	lc login.LoginServiceClient
}

func NewUserRpcDomain() *UserRpcDomain {
	return &UserRpcDomain{
		lc: rpc.LoginServiceClient,
	}
}

func (d *UserRpcDomain) MemberList(mIdList []int64) ([]*login.MemberMessage, map[int64]*login.MemberMessage, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	messageList, err := d.lc.FindMemInfoByIds(c, &login.UserMessage{MIds: mIdList})
	mMap := make(map[int64]*login.MemberMessage)
	for _, v := range messageList.List {
		mMap[v.Id] = v
	}
	return messageList.List, mMap, err
}

func (d *UserRpcDomain) MemberInfo(ctx context.Context, memberCode int64) (*login.MemberMessage, error) {
	c, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberMessage, err := d.lc.FindMemInfoById(c, &login.UserMessage{MemId: memberCode})
	return memberMessage, err
}
