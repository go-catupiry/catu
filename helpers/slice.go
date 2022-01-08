package helpers

func SliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// TODO! test this method result ..
func SliceRemove(s []string, r string) ([]string, bool) {
	var newList []string
	removed := false

	for _, v := range s {
		if v == r {
			removed = true
			continue
		}

		newList = append(newList, v)
	}

	return newList, removed
}
