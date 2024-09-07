package models

import (
	"time"
)

type User struct {
	ID                 string     `json:"id"`
	Nickname           string     `json:"nickname"`
	Email              string     `json:"email"`
	Password           string     `json:"password"`
	FirstName          *string    `json:"first_name"`
	SecondName         *string    `json:"second_name"`
	FirstSurname       *string    `json:"first_surname"`
	SecondSurname      *string    `json:"second_surname"`
	Age                *int32     `json:"age"`
	TypeDocument       *string    `json:"type_document"`
	DocumentNumber     string     `json:"document_number"`
	Cellphone          string     `json:"cellphone"`
	Gender             *string    `json:"gender"`
	Nationality        *string    `json:"nationality"`
	Country            *string    `json:"country"`
	Department         *string    `json:"department"`
	City               *string    `json:"city"`
	RealIp             string     `json:"real_ip"`
	StatusId           int32      `json:"status_id"`
	FailedAttempts     int32      `json:"failed_attempts"`
	BlockDate          *time.Time `json:"block_date"`
	DisabledDate       *time.Time `json:"disabled_date"`
	LastLogin          *time.Time `json:"last_login"`
	LastChangePassword *time.Time `json:"last_change_password"`
	BirthDate          *time.Time `json:"birth_date"`
	VerifiedCode       *string    `json:"verified_code"`
	IsDeleted          bool       `json:"is_deleted"`
	DeletedAt          *time.Time `json:"deleted_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
