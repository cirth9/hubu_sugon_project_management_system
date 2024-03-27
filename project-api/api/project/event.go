package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"strconv"
	"test.com/project-api/pkg/model/pro"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/project"
	"time"
)

type HandlerEvent struct {
}

func NewEvent() *HandlerEvent {
	return &HandlerEvent{}
}

func (h *HandlerEvent) myEventProjectList(c *gin.Context) {
	result := &common.Result{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	memberId := c.GetInt64("memberId")
	memberName := c.GetString("memberName")
	page, _ := strconv.Atoi(c.PostForm("page"))
	pageSize, _ := strconv.Atoi(c.PostForm("pageSize"))

	deleted := c.PostForm("deleted") == "1"
	msg := &project.ProjectRpcMessage{
		MemberId:   memberId,
		MemberName: memberName,
		Page:       int64(page),
		PageSize:   int64(pageSize),
		Deleted:    deleted,
	}
	myProjectResponse, err := ProjectServiceClient.FindProjectByMemId(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}

	var pms []*pro.ProjectAndMember
	copier.Copy(&pms, myProjectResponse.Pm)
	if pms == nil {
		pms = []*pro.ProjectAndMember{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"list":  pms, //null nil -> []
		"total": myProjectResponse.Total,
	}))
}
