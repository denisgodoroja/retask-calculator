package calculator

import (
	"fmt"
)

type Result struct {
	Packs      map[int]int
	PacksCount int
	ItemsCount int
}

func Calculate(amount int, packSizes []int) map[int]int {
	results := make(map[int]int) // {size: qty, ...}

	for i := len(packSizes) - 1; i >= 0; i-- {
		a := int(float64(amount) / float64(packSizes[i]))
		amount -= packSizes[i] * a
		if i == 0 && amount > 0 {
			a++
			amount = 0
		}
		if a > 0 {
			results[packSizes[i]] = a
		}
	}

	fmt.Println("orderPacks", amount, results)

	return results
}
