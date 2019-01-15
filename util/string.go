package util

import "bytes"

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func SplitStringByN(s string, n int) (split []string) {
	sub := ""
	runes := bytes.Runes([]byte(s))
	length := len(runes)
	for i, r := range runes {
		sub = sub + string(r)
		if (i+1)%n == 0 {
			split = append(split, sub)
			sub = ""
		} else if (i + 1) == length {
			split = append(split, sub)
		}
	}
	return
}
