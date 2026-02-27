package service

import "github.com/Ishkhan88/go-study/internal/apperr"

// Переэкспортируем ошибки, чтобы не ломать существующие хендлеры/вызовы.
var (
	ErrBadInput = apperr.ErrBadInput
	ErrNotFound = apperr.ErrNotFound
)
