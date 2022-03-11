package domain

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
	Id uuid.UUID `gorm:"type:uuid"`
	FirstName string
	LastName string
	Email string
	Password string
	FileUrl string
	CreatedAt time.Time
}	

// first_name -> Varchar,
// last_name -> Varchar,
// email -> Varchar,
// password -> Varchar,
// file_url -> Varchar,
// created_at -> Timestamp,