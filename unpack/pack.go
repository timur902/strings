package unpack

import "strings"

func Pack(s string) string {

	if len(s) == 0 {
		return ""
	}

	rs := []rune(s)

	var builder strings.Builder

	count := 1

	for i := 1; i <= len(rs); i++ {

		if i < len(rs) && rs[i] == rs[i-1] {
			count++
			continue
		}

		builder.WriteRune(rs[i-1])

		if count > 1 {
			builder.WriteString(string(rune('0' + count)))
		}

		count = 1
	}

	return builder.String()
}