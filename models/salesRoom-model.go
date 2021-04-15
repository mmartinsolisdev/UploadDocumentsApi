package models

type SalesRooms struct {
	SalesRoomID string `gorm:"column:srID"`
	ContractCode string `gorm:"column:srContractCode"`
}

/*type Fields struct {
	Id    string `gorm:"column:cxID"`
	DocType string `gorm:"column:cxDocType"`
	Language  string `gorm:"column:cxla"`
	DocName string `gorm:"column:cxDocName"`
	ContractCode int `gorm:"column:cxContractCode"`
}*/

func (SalesRooms) TableName() string {
    return "SalesRooms"
}
