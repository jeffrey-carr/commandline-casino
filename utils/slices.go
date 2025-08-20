package utils

import "math/rand"

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
