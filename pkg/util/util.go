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

// Utility for building maps with all elements of arr as keys and val as value.
func MapOf[Tk comparable, Tv any](arr []Tk, val Tv) map[Tk]Tv {
	out := map[Tk]Tv{}
	for _, e := range arr {
		out[e] = val
	}
	return out
}
