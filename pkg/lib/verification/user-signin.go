package verification

func LoginAndPassword(login, password string) bool {
	if login == "" || password == "" {
		return false
	}
	return true
}
