package user_basic_service_v1

import (
	"context"
	"github.com/jinzhu/copier"
	"test.com/project-grpc/user/user_basic"
	"test.com/project-user/internal/dao"
	"test.com/project-user/internal/data/member"
	"test.com/project-user/internal/repo"
)

type UserBasicService struct {
	user_basic.UnimplementedUserBasicServiceServer
	memberRepo repo.MemberRepo
}

func New() *UserBasicService {
	return &UserBasicService{
		memberRepo: dao.NewMemberDao(),
	}
}

func (u *UserBasicService) UpdateUserInfo(ctx context.Context, userBasicInfo *user_basic.MemberMessage) (*user_basic.MemberMessage, error) {
	var updateMemberInfo *member.Member
	copier.Copy(updateMemberInfo, userBasicInfo)
	_, err := u.memberRepo.UpdateMember(ctx, updateMemberInfo)
	if err != nil {
		return nil, err
	}
	return userBasicInfo, err
}
