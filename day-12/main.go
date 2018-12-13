package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type node struct {
	HasPlant   bool
	Next, Prev *node
}

type rule struct {
	Inputs        [5]bool
	ProducesPlant bool
}

func ensureListSize(head *node) *node {
	lowest, highest := head, head
	for walk := head; walk != nil; walk = walk.Prev {
		if walk.HasPlant {
			lowest = walk
		}
	}
	for walk := head; walk != nil; walk = walk.Next {
		if walk.HasPlant {
			highest = walk
		}
	}

	for i := 0; i < 2; i++ {
		if lowest.Prev == nil {
			newNode := new(node)
			newNode.Next = lowest
			lowest.Prev = newNode
		}
		lowest = lowest.Prev

		if highest.Next == nil {
			newNode := new(node)
			newNode.Prev = highest
			highest.Next = newNode
		}
		highest = highest.Next
	}

	return lowest
}

func loadData(filename string) (*node, []*rule) {
	file, _ := os.Open(filename)
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	line := scanner.Text()
	initialStateString := strings.Split(line, " ")[2]
	initialState := new(node)
	prev := initialState
	var cur *node
	for _, char := range initialStateString {
		cur = new(node)
		if char == '#' {
			cur.HasPlant = true
		}
		prev.Next = cur
		cur.Prev = prev
		prev = cur
	}
	initialState.Next.Prev = nil
	initialState = initialState.Next
	for i := 0; i < 2; i++ {
		newNode := new(node)
		newNode.Next = initialState
		initialState.Prev = newNode
		initialState = newNode
	}

	rules := make([]*rule, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		newRule := new(rule)
		for i := 0; i < 5; i++ {
			if line[i] == '#' {
				newRule.Inputs[i] = true
			}
		}
		if line[len(line)-1] == '#' {
			newRule.ProducesPlant = true
		}
		rules = append(rules, newRule)
	}

	file.Close()

	return initialState, rules
}

func testRule(plant *node, plantRule *rule) bool {
	if plant.HasPlant != plantRule.Inputs[2] {
		return false
	}

	if plant.Prev == nil && (plantRule.Inputs[0] != false || plantRule.Inputs[1] != false) {
		return false
	}
	if plant.Prev != nil {
		if plant.Prev.HasPlant != plantRule.Inputs[1] {
			return false
		}
		if plant.Prev.Prev == nil || plant.Prev.Prev.HasPlant != plantRule.Inputs[0] {
			return false
		}
	}

	if plant.Next == nil && (plantRule.Inputs[3] != false || plantRule.Inputs[4] != false) {
		return false
	}
	if plant.Next != nil {
		if plant.Next.HasPlant != plantRule.Inputs[3] {
			return false
		}
		if plant.Next.Next == nil || plant.Next.Next.HasPlant != plantRule.Inputs[4] {
			return false
		}
	}

	return true
}

func stepGeneration(head *node, rules []*rule) *node {
	head = ensureListSize(head)
	dict := make(map[*node]bool)
	for walk := head; walk != nil; walk = walk.Next {
		for i := range rules {
			if testRule(walk, rules[i]) {
				dict[walk] = rules[i].ProducesPlant
				break
			}
		}
	}
	for plant, newValue := range dict {
		plant.HasPlant = newValue
	}
	return head
}

func countPlants(head *node) int {
	sum := 0
	for walk := head; walk != nil; walk = walk.Next {
		if walk.HasPlant {
			sum++
		}
	}
	return sum
}

func plantsToString(head *node) string {
	var builder strings.Builder
	for walk := head; walk != nil; walk = walk.Next {
		if walk.HasPlant {
			builder.WriteString("#")
		} else {
			builder.WriteString(".")
		}
	}
	return builder.String()
}

func main() {
	head, rules := loadData("input.test.txt")
	for i := 0; i <= 20; i++ {
		fmt.Printf("%2d: %s\n", i, plantsToString(head))
		head = stepGeneration(head, rules)
	}
}
