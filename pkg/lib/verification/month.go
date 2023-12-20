package verification

func monthVer(month string) bool {
	if len(month) < 1 || len(month) > 2 {
		return false
	}
	if len(month) == 1 && month >= "1" && month <= "9" {
		return true
	}
	if len(month) == 2 && month >= "10" && month <= "12" {
		return true
	}
	return false
}
