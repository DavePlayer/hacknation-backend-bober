package models

type Organization struct {
	ID   string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
	// Relacja odwrotna do użytkowników
	Users []User `gorm:"foreignKey:OrganizationID" json:"users,omitempty"`
}
