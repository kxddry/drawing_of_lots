package main

import (
	"github.com/sgade/randomorg"
	"math/rand"
	"sort"
)

func shuffle(slice []int64) []int64 {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// this approach uses O(n) time and O(n) space, but is truly random.
func shuffleTrulyRandom(slice []int64) ([]int64, error) {
	n := len(slice)
	random := randomorg.NewRandom(randomToken)
	decimals, err := random.GenerateDecimalFractions(n, 10)
	if err != nil {
		// if, for some reason, we couldn't connect to the random.org api, use the regular shuffle instead.
		return shuffle(slice), err
	}

	// assign each random float to each value (userIds in our case)
	type pair struct {
		value int64
		rnd   float64
	}
	pairs := make([]pair, n)
	for i, v := range slice {
		pairs[i] = pair{value: v, rnd: decimals[i]}
	}

	// sort based on random floats
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].rnd < pairs[j].rnd
	})

	// Write the shuffled values back into the original slice.
	for i, p := range pairs {
		slice[i] = p.value
	}

	return slice, nil
}
