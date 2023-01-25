package dto

import "github.com/pkg/errors"

var (
	ErrEventNotFound           = errors.New("event not found")
	ErrStorageTypeIsNotCorrect = errors.New("storage type is not correct")
)
