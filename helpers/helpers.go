package helpers

// ArrayStringContains checks if a string is contained in a string array
func ArrayStringContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
