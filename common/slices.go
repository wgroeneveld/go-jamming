package common

func Includes(slice []string, elem string) bool {
	for _, el := range slice {
		if el == elem {
			return true
		}
	}
	return false
}
