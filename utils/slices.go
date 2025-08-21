package utils

import (
	"math/rand"
	"slices"
)

func Map[T any, K any](arr []T, f func(int, T) K) []K {
	var results []K
	for i, v := range arr {
		results = append(results, f(i, v))
	}

	return results
}

func Flatten[T any](arr [][]T) []T {
	var flattened []T
	for _, inner := range arr {
		flattened = append(flattened, inner...)
	}
	return flattened
}

func Shuffle[T any](arr []T) {
	for i := range arr {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func Pop[T any](arr []T) T {
	var ret T
	if len(arr) == 0 {
		return ret
	}

	ret = arr[0]
	arr = arr[1:]
	return ret
}

func Any[T any](arr []T, f func(T) bool) bool {
	return slices.ContainsFunc(arr, f)
}

func All[T any](arr []T, f func(T) bool) bool {
	fWrapper := func(i T) bool { return !f(i) }
	return Any(arr, fWrapper)
}

func Dedupe[T comparable](arr []T) []T {
	set := make(map[T]struct{}, len(arr))
	for _, i := range arr {
		set[i] = struct{}{}
	}

	remaining := make([]T, 0, len(set))
	for i := range set {
		remaining = append(remaining, i)
	}

	return remaining
}

// MaxFunc gets the max item from the slice given the ranking function
func MaxFunc[T any](arr []T, f func(T) int) (T, bool) {
	var maxItem T

	if len(arr) == 0 {
		return maxItem, false
	}

	maxValue := -1
	for _, i := range arr {
		v := f(i)
		if v > maxValue {
			maxItem = i
			maxValue = v
		}
	}

	return maxItem, true
}
