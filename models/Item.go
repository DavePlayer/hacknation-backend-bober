package models

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusReturned Status = "returned"
)

type Item struct {
	ID                     int64     `json:"id"`
	Issuer_id              string    `json:"issuer"`
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
	Status                 Status    `json:"status"`
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
	Status                 Status    `json:"status"`
}
