package tasks

type TaskMember struct {
	Id                int64  `json:"id"`
	Name              string `json:"name"`
	Avatar            string `json:"avatar"`
	Code              string `json:"code"`
	IsExecutor        int    `json:"is_executor"`
	IsOwner           int    `json:"is_owner"`
	MembarAccountCode string `json:"membar_account_code"`
}
