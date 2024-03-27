package data

type File struct {
	Id               int64
	PathName         string
	Title            string
	Extension        string
	Size             int
	ObjectType       string
	OrganizationCode int64
	TaskCode         int64
	ProjectCode      int64
	CreateBy         int64
	CreateTime       int64
	Downloads        int
	Extra            string
	Deleted          int
	FileUrl          string
	FileType         string
	DeletedTime      int64
}

func (*File) TableName() string {
	return "ms_file"
}
