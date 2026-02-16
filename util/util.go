package util

import "strings"

// Check if a string is really empty by trimming spaces.
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// Returns the values of map `m`.
func Values[Tk comparable, Tv any](m map[Tk]Tv) []Tv {
	a := []Tv{}
	for _, v := range m {
		a = append(a, v)
	}
	return a
}

// Returns a map whose keys are taken from `keys`,
// with every key initialized to `v`.
func MapWithValue[Tk comparable, Tv any](keys []Tk, v Tv) map[Tk]Tv {
	out := map[Tk]Tv{}
	for _, key := range keys {
		out[key] = v
	}
	return out
}
