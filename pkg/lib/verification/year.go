package verification

func Year(year string) bool {
	if len(year) != 4 {
		return false
	}
	if year[0] != '2' {
		return false
	}
	return true
}
