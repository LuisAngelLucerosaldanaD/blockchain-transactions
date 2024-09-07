package files

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// File  Model struct File
type File struct {
	ID           int64     `json:"id" db:"id" valid:"-"`
	IdDocument   int64     `json:"id_document" db:"id_document" valid:"required"`
	OriginalFile string    `json:"original_file" db:"original_file" valid:"required"`
	Hash         string    `json:"hash" db:"hash" valid:"required"`
	FileSize     int       `json:"file_size" db:"file_size" valid:"required"`
	Path         string    `json:"path" db:"path" valid:"required"`
	FileName     string    `json:"file_name" db:"file_name" valid:"required"`
	NumberPage   int       `json:"number_page" db:"number_page" valid:"required"`
	Bucket       string    `json:"bucket" db:"bucket" valid:"required"`
	IdFile       int       `json:"id_file" db:"id_file" valid:"required"`
	IdUser       int       `json:"id_user" db:"id_user" valid:"required"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Encoding     string    `json:"encoding"`
}

type ResponseFile struct {
	NameDocument string `json:"name_document"`
	Encoding     string `json:"encoding"`
	FileID       int    `json:"file_id"`
}

func NewUploadFile(idDocument int64, originalFile string, enconding string) *File {
	return &File{
		IdDocument:   idDocument,
		OriginalFile: originalFile,
		Encoding:     enconding,
	}
}

func (m *File) valid() (bool, error) {
	result, err := govalidator.ValidateStruct(m)
	if err != nil {
		return result, err
	}
	return result, nil
}
