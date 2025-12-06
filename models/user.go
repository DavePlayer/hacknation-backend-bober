package models

import "time"

type User struct {
	ID             string       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name           string       `gorm:"size:100;not null" json:"name"`
	Surname        string       `gorm:"size:100;not null" json:"surname"`
	OrganizationID string       `gorm:"not null" json:"organizationId"` // klucz obcy
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	CreatedAt      time.Time    `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time    `gorm:"autoUpdateTime" json:"updatedAt"`
}
