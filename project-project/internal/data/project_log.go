package data

import (
	"github.com/jinzhu/copier"
	"test.com/project-common/encrypts"
	"test.com/project-common/tms"
)

type ProjectLog struct {
	Id           int64
	MemberCode   int64
	Content      string
	Remark       string
	Type         string
	CreateTime   int64
	SourceCode   int64
	ActionType   string
	ToMemberCode int64
	IsComment    int
	ProjectCode  int64
	Icon         string
	IsRobot      int
}

func (*ProjectLog) TableName() string {
	return "ms_project_log"
}

type ProjectLogDisplay struct {
	Id           int64
	MemberCode   string
	Content      string
	Remark       string
	Type         string
	CreateTime   string
	SourceCode   string
	ActionType   string
	ToMemberCode string
	IsComment    int
	ProjectCode  string
	Icon         string
	IsRobot      int
	Member       Member
}

func (l *ProjectLog) ToDisplay() *ProjectLogDisplay {
	pd := &ProjectLogDisplay{}
	copier.Copy(pd, l)
	pd.MemberCode = encrypts.EncryptNoErr(l.MemberCode)
	pd.ToMemberCode = encrypts.EncryptNoErr(l.ToMemberCode)
	pd.ProjectCode = encrypts.EncryptNoErr(l.ProjectCode)
	pd.CreateTime = tms.FormatByMill(l.CreateTime)
	pd.SourceCode = encrypts.EncryptNoErr(l.SourceCode)
	return pd
}

type IndexProjectLogDisplay struct {
	Content      string
	Remark       string
	CreateTime   string
	SourceCode   string
	IsComment    int
	ProjectCode  string
	MemberAvatar string
	MemberName   string
	ProjectName  string
	TaskName     string
}

func (l *ProjectLog) ToIndexDisplay() *IndexProjectLogDisplay {
	pd := &IndexProjectLogDisplay{}
	copier.Copy(pd, l)
	pd.ProjectCode = encrypts.EncryptNoErr(l.ProjectCode)
	pd.CreateTime = tms.FormatByMill(l.CreateTime)
	pd.SourceCode = encrypts.EncryptNoErr(l.SourceCode)
	return pd
}
