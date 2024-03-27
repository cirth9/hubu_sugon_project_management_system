package data

import (
	"github.com/jinzhu/copier"
	"test.com/project-common/encrypts"
	"test.com/project-common/tms"
)

type SourceLink struct {
	Id               int64
	SourceType       string
	SourceCode       int64
	LinkType         string
	LinkCode         int64
	OrganizationCode int64
	CreateBy         int64
	CreateTime       int64
	Sort             int
}

func (*SourceLink) TableName() string {
	return "ms_source_link"
}

type SourceLinkDisplay struct {
	Id               int64        `json:"id"`
	Code             string       `json:"code"`
	SourceType       string       `json:"source_type"`
	SourceCode       string       `json:"source_code"`
	LinkType         string       `json:"link_type"`
	LinkCode         string       `json:"link_code"`
	OrganizationCode string       `json:"organization_code"`
	CreateBy         string       `json:"create_by"`
	CreateTime       string       `json:"create_time"`
	Sort             int          `json:"sort"`
	Title            string       `json:"title"`
	SourceDetail     SourceDetail `json:"sourceDetail"`
}

type SourceDetail struct {
	Id               int64  `json:"id"`
	Code             string `json:"code"`
	PathName         string `json:"path_name"`
	Title            string `json:"title"`
	Extension        string `json:"extension"`
	Size             int    `json:"size"`
	ObjectType       string `json:"object_type"`
	OrganizationCode string `json:"organization_code"`
	TaskCode         string `json:"task_code"`
	ProjectCode      string `json:"project_code"`
	CreateBy         string `json:"create_by"`
	CreateTime       string `json:"create_time"`
	Downloads        int    `json:"downloads"`
	Extra            string `json:"extra"`
	Deleted          int    `json:"deleted"`
	FileUrl          string `json:"file_url"`
	FileType         string `json:"file_type"`
	DeletedTime      string `json:"deleted_time"`
	ProjectName      string `json:"projectName"`
	FullName         string `json:"fullName"`
}

func (s *SourceLink) ToDisplay(f *File) *SourceLinkDisplay {
	sl := &SourceLinkDisplay{}
	copier.Copy(sl, s)
	sl.SourceDetail = SourceDetail{}
	copier.Copy(&sl.SourceDetail, f)
	sl.LinkCode = encrypts.EncryptNoErr(s.LinkCode)
	sl.OrganizationCode = encrypts.EncryptNoErr(s.OrganizationCode)
	sl.CreateTime = tms.FormatByMill(s.CreateTime)
	sl.CreateBy = encrypts.EncryptNoErr(s.CreateBy)
	sl.SourceCode = encrypts.EncryptNoErr(s.SourceCode)
	sl.SourceDetail.OrganizationCode = encrypts.EncryptNoErr(f.OrganizationCode)
	sl.SourceDetail.CreateBy = encrypts.EncryptNoErr(f.CreateBy)
	sl.SourceDetail.CreateTime = tms.FormatByMill(f.CreateTime)
	sl.SourceDetail.DeletedTime = tms.FormatByMill(f.DeletedTime)
	sl.SourceDetail.TaskCode = encrypts.EncryptNoErr(f.TaskCode)
	sl.SourceDetail.ProjectCode = encrypts.EncryptNoErr(f.ProjectCode)
	sl.SourceDetail.FullName = f.Title
	sl.Title = f.Title
	return sl
}
