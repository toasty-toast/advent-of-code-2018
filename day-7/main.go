package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type step struct {
	Value        string
	Dependencies []*step
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadSteps(filename string) []*step {
	file, err := os.Open(filename)
	check(err)

	existingSteps := make(map[string]*step)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		value := strings.Split(line, " ")[7]
		dependency := strings.Split(line, " ")[1]

		var dependentStep *step
		if _, contains := existingSteps[value]; !contains {
			dependentStep = new(step)
			dependentStep.Value = value
			dependentStep.Dependencies = make([]*step, 0)
			existingSteps[value] = dependentStep
		} else {
			dependentStep = existingSteps[value]
		}

		var dependencyStep *step
		if _, contains := existingSteps[dependency]; !contains {
			dependencyStep = new(step)
			dependencyStep.Value = dependency
			dependencyStep.Dependencies = make([]*step, 0)
			existingSteps[dependency] = dependencyStep
		} else {
			dependencyStep = existingSteps[dependency]
		}

		dependentStep.Dependencies = append(dependentStep.Dependencies, dependencyStep)
	}

	steps := make([]*step, 0)
	for _, value := range existingSteps {
		steps = append(steps, value)
	}

	file.Close()
	return steps
}

func orderSteps(steps []*step) string {
	done := make(map[*step]bool)
	availableSteps := func() []*step {
		available := make([]*step, 0)
		for _, value := range steps {
			if _, ok := done[value]; ok {
				continue
			}
			depsSatisfied := true
			for _, dependency := range value.Dependencies {
				if _, ok := done[dependency]; !ok {
					depsSatisfied = false
					break
				}
			}
			if depsSatisfied {
				available = append(available, value)
			}
		}
		return available
	}

	var builder strings.Builder
	for available := availableSteps(); len(available) > 0; available = availableSteps() {
		sort.Slice(available, func(i, j int) bool {
			return available[i].Value < available[j].Value
		})
		builder.WriteString(available[0].Value)
		done[available[0]] = true
	}
	return builder.String()
}

func timeSteps(steps []*step, numWorkers int) int {
	done := make(map[*step]bool)
	availableSteps := func() []*step {
		available := make([]*step, 0)
		for _, value := range steps {
			if _, ok := done[value]; ok {
				continue
			}
			depsSatisfied := true
			for _, dependency := range value.Dependencies {
				if _, ok := done[dependency]; !ok {
					depsSatisfied = false
					break
				}
			}
			if depsSatisfied {
				available = append(available, value)
			}
		}
		return available
	}

	timer := 0
	assigned := make(map[*step]bool)
	assignments := make([]*step, numWorkers)
	assignmentTimers := make([]int, numWorkers)
	for {
		for i := range assignments {
			if assignments[i] != nil && int(assignments[i].Value[0]-'A'+61) == assignmentTimers[i] {
				done[assignments[i]] = true
				assignments[i] = nil
				assignmentTimers[i] = 0
				delete(assigned, assignments[i])
			}
		}

		available := availableSteps()
		sort.Slice(available, func(i, j int) bool {
			return available[i].Value < available[j].Value
		})
		if len(available) == 0 {
			break
		}

		for i := range available {
			for j := range assignments {
				if assignments[j] == nil {
					if _, ok := assigned[available[i]]; !ok {
						assignments[j] = available[i]
						assigned[assignments[j]] = true
						break
					}
				}
			}
		}

		for i := range assignments {
			if assignments[i] != nil {
				assignmentTimers[i]++
			}
		}
		timer++
	}

	return timer
}

func main() {
	steps := loadSteps("input.txt")
	fmt.Printf("Instruction order: %s\n", orderSteps(steps))
	fmt.Printf("All steps took %ds\n", timeSteps(steps, 5))
}
