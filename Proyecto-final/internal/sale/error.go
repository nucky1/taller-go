package sale

import "errors"

var (
	ErrSaleNotFound       = errors.New("sale not found")
	ErrInvalidStateChange = errors.New("transición de estado no permitida")
	ErrInvalidNewState    = errors.New("estado no válido para cambio")
)
