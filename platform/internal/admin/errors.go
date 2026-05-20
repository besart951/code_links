package admin

import "errors"

var (
	ErrSuperadminRequired = errors.New("superadmin required")
	ErrNotFound           = errors.New("admin record not found")
)
