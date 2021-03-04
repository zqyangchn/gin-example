package app

func GetPageOffset(pageNumber, pageSize int) int {
	if pageNumber > 0 {
		return (pageNumber - 1) * pageSize
	}

	return 0
}
