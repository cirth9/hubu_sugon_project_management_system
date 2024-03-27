package project

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"net/http"
	"test.com/project-api/pkg/model"
	common "test.com/project-common"
	"test.com/project-common/errs"
	"test.com/project-grpc/department"
	"time"
)

type HandlerDepartment struct {
}

func (d *HandlerDepartment) department(c *gin.Context) {
	result := &common.Result{}
	var req *model.DepartmentReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &department.DepartmentReqMessage{
		Page:                 req.Page,
		PageSize:             req.PageSize,
		ParentDepartmentCode: req.Pcode,
		OrganizationCode:     c.GetString("organizationCode"),
	}
	listDepartmentMessage, err := DepartmentServiceClient.List(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var list []*model.Department
	copier.Copy(&list, listDepartmentMessage.List)
	if list == nil {
		list = []*model.Department{}
	}
	c.JSON(http.StatusOK, result.Success(gin.H{
		"total": listDepartmentMessage.Total,
		"page":  req.Page,
		"list":  list,
	}))
}

func (d *HandlerDepartment) save(c *gin.Context) {
	result := &common.Result{}
	var req *model.DepartmentReq
	c.ShouldBind(&req)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &department.DepartmentReqMessage{
		Name:                 req.Name,
		DepartmentCode:       req.DepartmentCode,
		ParentDepartmentCode: req.ParentDepartmentCode,
		OrganizationCode:     c.GetString("organizationCode"),
	}
	departmentMessage, err := DepartmentServiceClient.Save(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var res = &model.Department{}
	copier.Copy(res, departmentMessage)
	c.JSON(http.StatusOK, result.Success(res))
}

func (d *HandlerDepartment) read(c *gin.Context) {
	result := &common.Result{}
	departmentCode := c.PostForm("departmentCode")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	msg := &department.DepartmentReqMessage{
		DepartmentCode:   departmentCode,
		OrganizationCode: c.GetString("organizationCode"),
	}
	departmentMessage, err := DepartmentServiceClient.Read(ctx, msg)
	if err != nil {
		code, msg := errs.ParseGrpcError(err)
		c.JSON(http.StatusOK, result.Fail(code, msg))
	}
	var res = &model.Department{}
	copier.Copy(res, departmentMessage)
	c.JSON(http.StatusOK, result.Success(res))
}

func NewDepartment() *HandlerDepartment {
	return &HandlerDepartment{}
}
