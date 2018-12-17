package main

import "fmt"

const (
	puzzleInput       = 409551
	fistRecipeValue   = 3
	secondRecipeValue = 7
)

func newRecipes(inputs []int) []int {
	sum := 0
	fmt.Print("Inputs: ")
	for i := range inputs {
		fmt.Printf("%d ", inputs[i])
		sum += inputs[i]
	}

	results := make([]int, 0)
	for sum > 0 {
		results = append(results, sum%10)
		sum /= 10
	}
	for i := 0; i < len(results)/2; i++ {
		temp := results[i]
		results[i] = results[len(results)-1-i]
		results[len(results)-i-1] = temp
	}

	fmt.Print("Outputs: ")
	for i := range results {
		fmt.Printf("%d ", results[i])
	}
	fmt.Println()

	return results
}

func printScores(scores []int, elf1 int, elf2 int) {
	for i := range scores {
		if i == elf1 {
			fmt.Printf("(%d) ", scores[i])

		} else if i == elf2 {
			fmt.Printf("[%d] ", scores[i])

		} else {
			fmt.Printf("%d ", scores[i])
		}
	}
	fmt.Println()
}

func main() {
	scores := []int{fistRecipeValue, secondRecipeValue}
	elf1, elf2 := 0, 1
	input := 2018
	value := (input * 2) - 2
	//printScores(scores, elf1, elf2)
	for i := 0; i < value; i++ {
		newRecipes := newRecipes([]int{scores[elf1], scores[elf2]})
		for _, recipe := range newRecipes {
			//fmt.Printf("Adding %d\n", recipe)
			scores = append(scores, recipe)
		}

		elf1 = (elf1 + scores[elf1] + 1) % len(scores)
		elf2 = (elf2 + scores[elf2] + 1) % len(scores)

		//printScores(scores, elf1, elf2)
	}

	fmt.Print("Scores: ")
	for i := 0; i < 10; i++ {
		fmt.Print(scores[input+i])
	}
	fmt.Println()
}
