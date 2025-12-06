package models

type User struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Surname      string       `json:"surname"`
	Organization Organization `json:"organization"`
}
