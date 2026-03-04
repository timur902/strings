package repository

import "github.com/google/uuid"

type Repository interface {
	InsertString(str string) (uuid.UUID, error) //uuid, error
	SelectByID(id uuid.UUID) (string, error) // найденная строка/ не найдено 
}
