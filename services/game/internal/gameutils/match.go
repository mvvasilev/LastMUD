package gameutils

import "slices"

func OneOf[T comparable](value T, tests ...T) bool {
	return slices.Contains(tests, value)
}
