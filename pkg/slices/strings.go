package slices

func AddToSet(slice []string, items ...string) []string {
	for _, item := range items {
		if Contain(slice, item) {
			continue
		}

		slice = append(slice, item)
	}

	return slice
}

func Contain(slice []string, item string) bool {
	for _, elmt := range slice {
		if elmt == item {
			return true
		}
	}

	return false
}

func Concat(a, b []string) []string {
	result := make([]string, 0, len(a)+len(b))

	for _, v := range a {
		if Contain(result, v) {
			continue
		}

		result = append(result, v)
	}

	for _, v := range b {
		if Contain(result, v) {
			continue
		}

		result = append(result, v)
	}

	return result
}

func Index(slice []string, item string) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == item {
			return i
		}
	}

	return -1
}

func Remove(target string, slice []string) []string {
	index := Index(slice, target)
	if index < 0 {
		return slice
	}

	resp := make([]string, 0, len(slice)-1)
	resp = append(resp, slice[:index]...)
	resp = append(resp, slice[index+1:]...)

	return resp
}
