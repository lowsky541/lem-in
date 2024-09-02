package lemin

func ArrayEquals[t comparable](a []t, b []t) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// Basically, this function takes a map and returns the values
// discarding the keys.
// Here, kt is the type of the key and vt the type of the values, so: map[kt]vt
// See TestMapValues() in utilities/generics_test.go.
func MapValues[kt comparable, vt any](m map[kt]vt) []vt {
	a := []vt{}
	for _, v := range m {
		a = append(a, v)
	}
	return a
}

// Basically, this function takes any element of type T and an array of elements
// of the same type T. It retuns whether or not the element is in the array.
// The type T must be comparable by == (string, int but not struct).
// This kind of function is a generic and can accept a wide range of types.
func Contains[T comparable](element T, array []T) bool {
	for _, v := range array {
		if v == element {
			return true
		}
	}
	return false
}

func ArrayRemove[T comparable](a []T, i int) []T {
	return append(a[:i], a[i+1:]...)
}
