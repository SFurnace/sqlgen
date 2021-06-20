package tests

type Customer struct {
	Uin               int64  `db:"uin"`
	AppID             int64  `db:"appId"`
	CustomerName      string `db:"userName"`
	CustomerIndustry  string `db:"userIndustry"`
	CustomerArchitect string `db:"userArchitect"`
	CustomerSeller    string `db:"userSeller"`
	RemarkName        string `db:"remarkName"`
	PicUrl            string `db:"picUrl"`        // 客户对应的url
	IndustryGrade     string `db:"industryGrade"` // 客户行业定级
}
