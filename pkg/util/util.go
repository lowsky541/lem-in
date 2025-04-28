package util

import "strings"

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func OnlyValues[Tk comparable, Tv any](m map[Tk]Tv) []Tv {
	a := []Tv{}
	for _, v := range m {
		a = append(a, v)
	}
	return a
}
