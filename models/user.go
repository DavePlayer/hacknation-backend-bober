package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string    `gorm:"size:100;not null" json:"email"`
	Password     string    `gorm:"size:100;not null" json:"password"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Surname      string    `gorm:"size:100;not null" json:"surname"`
	Organization string    `gorm:"size:100;not null"  json:"organization"`
	City         string    `gorm:"size:100; null"  json:"city"`
	Voivodeship  string    `gorm:"size:100; null"  json:"voivodeship"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type ReturnedUser struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string    `gorm:"size:100;not null" json:"email"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Surname      string    `gorm:"size:100;not null" json:"surname"`
	Organization string    `gorm:"size:100;not null"  json:"organization"`
	City         string    `gorm:"size:100; null"  json:"city"`
	Voivodeship  string    `gorm:"size:100; null"  json:"voivodeship"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ReturnedUser) From(u User) ReturnedUser {
	return ReturnedUser{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		Surname:      u.Surname,
		Organization: u.Organization,
		City:         u.City,
		Voivodeship:  u.Voivodeship,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
