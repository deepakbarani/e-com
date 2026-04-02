package dto

import "github.com/gofrs/uuid"

type Pagination struct {
	Page   int `json:"page" form:"page"`
	Offset int `json:"-" form:"-"`
	Limit  int `json:"-" form:"-"`
}

type ProductFilter struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Quantity int64     `json:"quantity"`
	AddedBy  uuid.UUID `json:"added_by"`
}
