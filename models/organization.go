package models

type Organization struct {
	ID    uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"size:255;not null" json:"name"`
	Users []User `gorm:"foreignKey:OrganizationID" json:"users,omitempty"`
}
