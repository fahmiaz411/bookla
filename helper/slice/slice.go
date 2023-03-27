package slice

func Includes[T comparable](lists []T, item T) bool {
	for _, v := range lists {
		if v == item {
			return true
		}
	}

	return false
}
