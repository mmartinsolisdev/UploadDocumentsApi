package models

type SalesRooms struct {
	SalesRoomID string `gorm:"column:srID"`
	ContractCode string `gorm:"column:srContractCode"`
}

func (SalesRooms) TableName() string {
    return "SalesRooms"
}
