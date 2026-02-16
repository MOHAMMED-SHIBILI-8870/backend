package models

type Category struct{
	ID uint `gorm:"PRIMARYKEY" json:"id"`
	CategoryName string `gorm:"size:50;not null;unique" json:"category_name"`
	ParentID *uint `gorm:"constraint:OnDelete:CASECADE" json:"parent_id"`
}