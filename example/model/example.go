package model

type Example struct {
	Id       uint       `json:"id" gorm:"column:id;`
	Name     NullString `json:"name" gorm:"column:name;`
	Phone    NullString `json:"phone" gorm:"type:char(11);column:phone;`
	Password NullString `json:"password" gorm:"column:password;`
}

func (m Example) TableName() string {
	return "example"
}
