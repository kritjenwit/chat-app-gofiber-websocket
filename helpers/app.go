package helpers

import "strconv"

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	r := make([]V, 0, len(m))
	for _, v := range m {
		r = append(r, v)
	}
	return r
}

func ToInt(s string) int {
	intVar, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return intVar
}
