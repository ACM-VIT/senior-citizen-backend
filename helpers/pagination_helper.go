package helpers

func CalculatePagination(page, pageSize int, total int64) (offset int, totalPages int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset = (page - 1) * pageSize
	totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	return
}
