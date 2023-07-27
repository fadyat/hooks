package helpers

func Find[T any](content []*T, predicate func(*T) bool) *T {
	for _, c := range content {
		if predicate(c) {
			return c
		}
	}

	return nil
}

func RemoveDuplicates[T comparable](content []*T) []*T {
	keys := make(map[T]bool)
	list := make([]*T, 0, len(content))

	for _, entry := range content {
		if _, value := keys[*entry]; !value {
			keys[*entry] = true
			list = append(list, entry)
		}
	}

	return list
}
