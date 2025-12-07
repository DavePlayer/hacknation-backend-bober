package models

import "time"

type Item struct {
	ID                     uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Issuer_id              uint      `gorm:"not null"  json:"issuerId"`
	Name                   string    `gorm:"size:100;not null" json:"itemName"`
	Type                   string    `gorm:"size:100;not null" json:"type"`
	Description            string    `gorm:"size:500;not null" json:"description"`
	Document_transfer_date time.Time `json:"documentTransferDate"` // time when document was given
	Entry_date             time.Time `gorm:"autoCreateTime" json:"entryDate"`
	Found_date             time.Time `json:"foundDate"`
	Issue_number           string    `gorm:"size:100;not null" json:"issueNumber"`
	Where_stored           string    `gorm:"size:100;not null" json:"whereStored"` // where it is stored
	Where_found            string    `gorm:"size:100;not null" json:"whereFound"`  // where it was found
	Voivodeship            string    `gorm:"size:100;not null" json:"voivodeship"` // where it is stored
	Status                 string    `gorm:"size:100;not null" json:"status"`
}

type ImportedItem struct {
	Name                   string    `json:"itemName"`
	Type                   string    `json:"type"`
	Description            string    `json:"description"`
	Document_transfer_date time.Time `json:"documentTransferDate"` // time when document was given
	Entry_date             time.Time `json:"entryDate"`
	Found_date             time.Time `json:"foundDate"`
	Issue_number           string    `json:"issueNumber"`
	Where_stored           string    `json:"whereStored"` // where it is stored
	Where_found            string    `json:"whereFound"`  // where it was found
	Voivodeship            string    `json:"voivodeship"` // where it is stored
	Status                 string    `json:"status"`
}
