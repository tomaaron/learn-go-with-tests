package iteration

func Repeat(s string, r int) (result string) {
	for i := 0; i < r; i++ {
		result += s
	}
	return
}
