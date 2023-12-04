package storage

import "errors"

var (
	ErrDBConnect    = errors.New("error on db connect")
	ErrDBClose      = errors.New("error on db close")
	ErrNotFound     = errors.New("entity not found")
	ErrAlreadyExist = errors.New("entity already exist")
)
