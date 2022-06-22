package common

func Shorten(txt string) string {
	if len(txt) <= 250 {
		return txt
	}
	return txt[:250] + "..."
}
