package util

import "iter"

// Hash-based intersection of slices
func IntersectSlices[T comparable](a []T, b []T) []T {
	set := make([]T, 0)
	hash := make(map[T]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	for _, v := range b {
		if _, ok := hash[v]; ok {
			set = append(set, v)
		}
	}

	return set
}

// Hash-based intersection of iterators
func IntersectIterators[T comparable](a, b iter.Seq[T]) []T {
	set := make([]T, 0)
	hash := make(map[T]struct{})

	for v := range a {
		hash[v] = struct{}{}
	}

	for v := range b {
		if _, ok := hash[v]; ok {
			set = append(set, v)
		}
	}

	return set
}

// Hash-based intersection of iterator and slice
func IntersectSliceWithIterator[T comparable](a []T, b iter.Seq[T]) []T {
	set := make([]T, 0)
	hash := make(map[T]struct{})

	for _, v := range a {
		hash[v] = struct{}{}
	}

	for v := range b {
		if _, ok := hash[v]; ok {
			set = append(set, v)
		}
	}

	return set
}
