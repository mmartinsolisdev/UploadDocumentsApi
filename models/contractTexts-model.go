package models

type ContractTexts struct {
	Id    string `gorm:"column:cxID"`
	Language  string `gorm:"column:cxla"`
	DocType string `gorm:"column:cxDocType"`
	DocName string `gorm:"column:cxDocName"`
	//DocFile string //`gorm:"column:cxTextBinary"`
	ContractCode int `gorm:"column:cxContractCode"`
}

/*type Fields struct {
	Id    string `gorm:"column:cxID"`
	DocType string `gorm:"column:cxDocType"`
	Language  string `gorm:"column:cxla"`
	DocName string `gorm:"column:cxDocName"`
	ContractCode int `gorm:"column:cxContractCode"`
}*/

func (ContractTexts) TableName() string {
    return "ContractTexts"
}
