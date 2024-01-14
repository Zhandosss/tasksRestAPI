package verification

import "strconv"

func Date(day, month, year string) bool {
	var (
		thirtyDayMonth = map[string]struct{}{
			"4":  {},
			"6":  {},
			"9":  {},
			"11": {},
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
		yearInt, _ := strconv.Atoi(year)
		if yearInt%4 == 0 && yearInt%100 != 0 || yearInt%400 == 0 {
			if day > "29" {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
