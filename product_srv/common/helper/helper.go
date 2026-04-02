package helper

import (
	"fmt"
	"product_srv/common/dto"

	"github.com/gofrs/uuid"
)

func StringToUUID(idStr string) (uuid.UUID, error) {

	id, err := uuid.FromString(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Failed to parsed string to UUID : %v", err)
	}
	return id, nil
}

func FindOffset(page, limit int) dto.Pagination {

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit
	pageParams := dto.Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
	return pageParams
}
