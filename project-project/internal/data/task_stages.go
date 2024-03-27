package data

type TaskStages struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ProjectCode int64  `json:"project_code"`
	Sort        int    `json:"sort"`
	Description string `json:"description"`
	CreateTime  int64  `json:"create_time"`
	Deleted     int    `json:"deleted"`
}

func (*TaskStages) TableName() string {
	return "ms_task_stages"
}

func ToTaskStagesMap(tsList []*TaskStages) map[int]*TaskStages {
	m := make(map[int]*TaskStages)
	for _, v := range tsList {
		m[v.Id] = v
	}
	return m
}
