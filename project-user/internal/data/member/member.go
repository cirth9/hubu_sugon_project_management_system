package member

type Member struct {
	Id              int64
	Account         string
	Password        string
	Name            string
	Mobile          string
	Realname        string
	CreateTime      int64
	Status          int
	LastLoginTime   int64
	Sex             int
	Avatar          string
	Idcard          string
	Province        int
	City            int
	Area            int
	Address         string
	Description     string
	Email           string
	DingtalkOpenid  string
	DingtalkUnionid string
	DingtalkUserid  string
}

func (*Member) TableName() string {
	return "ms_member"
}
