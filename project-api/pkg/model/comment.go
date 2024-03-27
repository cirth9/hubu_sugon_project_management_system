package model

type CommentReq struct {
	TaskCode string   `form:"taskCode"`
	Comment  string   `form:"comment"`
	Mentions []string `form:"mentions"`
}
