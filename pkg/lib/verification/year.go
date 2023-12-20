package verification

func yearVer(year string) bool {
	if len(year) != 4 {
		return false
	}
	if year[0] != '2' {
		return false
	}
	for i := range year {
		if year[i] < '0' || year[i] > '9' {
			return false
		}
	}
	return true
}
