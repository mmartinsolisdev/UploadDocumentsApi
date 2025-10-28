package models

type ContractTexts struct {
	Id    string `gorm:"column:cxID"`
	Language  string `gorm:"column:cxla"`
	DocType string `gorm:"column:cxDocType"`
	DocName string `gorm:"column:cxDocName"`
	SaleType string `gorm:"column:cxSaleType"`
	ContractCode int `gorm:"column:cxContractCode"`
}

func (ContractTexts) TableName() string {
    return "ContractTexts"
}
