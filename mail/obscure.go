package mail

import (
	"fmt"
	"strings"
)

// ObscureAddress obscures 'addr' by replacing all but the first and last characters in its major parts with '***' characters.
func ObscureAddress(addr string) string {

	parts := strings.Split(addr, "@")

	name := parts[0]
	dest := parts[1]

	name = strings.Replace(name, ".", "", -1)
	name = strings.Replace(name, "-", "", -1)

	name = obscure(name, false)
	dest = obscure(dest, true)

	return fmt.Sprintf("%s@%s", name, dest)
}

func obscure(raw string, preserve_last bool) string {

	parts := strings.Split(raw, ".")
	count := len(parts)

	obscured := make([]string, 0)

	max := count

	if preserve_last {
		max = count - 1
	}

	for i := 0; i < max; i++ {

		str := parts[i]

		runes := strings.Split(str, "")
		count_r := len(runes)

		new_str := "***"

		if count_r >= 3 {
			new_str = fmt.Sprintf("%s***%s", runes[0], runes[count_r-1])
		}

		obscured = append(obscured, new_str)
	}

	if preserve_last {
		obscured = append(obscured, parts[count-1])
	}

	return strings.Join(obscured, ".")
}
