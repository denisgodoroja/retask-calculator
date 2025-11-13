package calculator

import (
	"sort"
)

// represents a potential solution.
type result struct {
	packs    map[int]int
	totalSum int
	numPacks int
}

// returns true if 'r' is a better solution than 'other'
// based on the constraints: min excess first, then min packs.
func (r result) isBetterThan(other result, target int) bool {
	// If other has no packs and sum is 0 (empty or initial result), then r is better (if r is valid)
	if other.totalSum == 0 && other.numPacks == 0 {
		return true
	}

	rExcess := r.totalSum - target
	otherExcess := other.totalSum - target

	// Priority 1: Minimal Excess
	if rExcess != otherExcess {
		return rExcess < otherExcess
	}

	// Priority 2: Minimal Number of Packs
	return r.numPacks < other.numPacks
}

func Calculate(amount int, packSizes []int) map[int]int {
	if amount <= 0 || len(packSizes) == 0 {
		return map[int]int{}
	}

	// Sort pack sizes descending for better efficiency in recursion
	sortedSizes := make([]int, len(packSizes))
	copy(sortedSizes, packSizes)
	sort.Sort(sort.Reverse(sort.IntSlice(sortedSizes)))

	// Memoization cache to store optimal results for remaining amounts
	memo := make(map[int]result)

	// Optimization: If the amount is significantly larger than the largest pack,
	// we can pre-fill some packs of the smallest size to reduce recursion depth.
	prefill := 0
	if amount >= 100*sortedSizes[0] {
		prefill = int(amount/sortedSizes[0]) - 20
		amount -= prefill * sortedSizes[0]
	}

	finalRes := solve(amount, sortedSizes, memo)
	if prefill > 0 {
		finalRes.packs[sortedSizes[0]] += prefill
	}

	return finalRes.packs
}

// recursively finds the best combination for the target amount.
func solve(target int, packSizes []int, memo map[int]result) result {
	// Validate the target amount
	if target <= 0 {
		return result{packs: map[int]int{}, totalSum: 0, numPacks: 0}
	}

	// Check memoization cache
	if res, ok := memo[target]; ok {
		return res
	}

	var bestRes result

	// Try every pack size
	for _, p := range packSizes {
		// 1. Recurse: see what happens if we use this pack
		currentRes := solve(target-p, packSizes, memo)

		// 2. Construct the new result from this recursion
		newPacks := make(map[int]int)
		for k, v := range currentRes.packs {
			newPacks[k] = v
		}
		newPacks[p]++

		candidate := result{
			packs:    newPacks,
			totalSum: currentRes.totalSum + p,
			numPacks: currentRes.numPacks + 1,
		}

		// 3. Compare with the best result found so far for this specific target
		if candidate.isBetterThan(bestRes, target) {
			bestRes = candidate
		}

		// Optimization: If we found a perfect match (excess 0) with 1 pack,
		// it's impossible to beat for this specific recursion level, so break early.
		if bestRes.totalSum == target && bestRes.numPacks == 1 {
			break
		}
	}

	memo[target] = bestRes

	return bestRes
}
