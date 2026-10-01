package storage

import "errors"



var(
	ErrURLExists = errors.New("url exists")
	ErrUrlNotFound = errors.New("url not found")
)