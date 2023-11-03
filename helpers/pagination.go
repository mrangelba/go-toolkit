package helpers

import "strconv"

func ConvertLimitAndOffset(page int, size int) (int, int) {
	offset := -1
	limit := -1

	if page > 0 {
		offset = (page - 1) * size
	}

	if size > 0 {
		limit = size
	}

	return limit, offset
}

func ConvertPageAndSize(p string, s string) (int, int) {
	page, err := strconv.Atoi(p)

	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(s)

	if err != nil {
		limit = 20
	}

	return page, limit
}
