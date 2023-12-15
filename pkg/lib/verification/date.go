package verification

func date(day, month, year string) bool {
	var (
		thirtyDayMonth = map[string]struct{}{
			"4":  struct{}{},
			"6":  struct{}{},
			"9":  struct{}{},
			"11": struct{}{},
		}
	)
	if !monthVer(month) || !yearVer(year) {
		return false
	}
	if len(day) == 1 && day >= "1" && day <= "9" {
		return true
	}
	if len(day) != 2 || day < "10" || day > "31" {
		return false
	}
	if day >= "10" && day <= "28" {
		return true
	}
	if _, ok := thirtyDayMonth[month]; ok && day == "31" {
		return false
	}
	if month == "2" {
		yearInt, _ :=
	}
	return true
}
