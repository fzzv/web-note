package iteration

func Repeat(character string, times int) string {
	var repeated string
	for range times {
		repeated += character
	}
	return repeated
}
