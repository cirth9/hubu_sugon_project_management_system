package tasks

type TaskStagesResp struct {
	Name         string `json:"name"`
	ProjectCode  string `json:"project_code"`
	Sort         int    `json:"sort"`
	Description  string `json:"description"`
	CreateTime   string `json:"create_time"`
	Code         string `json:"code"`
	Deleted      int    `json:"deleted"`
	TasksLoading bool   `json:"tasksLoading"`
	FixedCreator bool   `json:"fixedCreator"`
	ShowTaskCard bool   `json:"showTaskCard"`
	Tasks        []int  `json:"tasks"`
	DoneTasks    []int  `json:"doneTasks"`
	UnDoneTasks  []int  `json:"unDoneTasks"`
}
