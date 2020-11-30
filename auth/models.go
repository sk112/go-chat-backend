package auth

import (
	// Used in formating
	_ "encoding/json"
	// Used in formating
	_ "gorm.io/gorm"
)

// User ...
type User struct {
	UserID    string `json:"userid" gorm:"primaryKey"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Firstname string
	Lastname  string
}
