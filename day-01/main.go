package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getInputValues(filename string) []int {
	file, err := os.Open("input.txt")
	check(err)
	
	values := make([]int, 0)
	fscanner := bufio.NewScanner(file)
    for fscanner.Scan() {
		line := fscanner.Text()
		lineValue, err := strconv.ParseInt(line, 10, 32)
		check(err)
		values = append(values, int(lineValue))
	}

	return values
}

func getResultingFrequency(values []int) int {
	sum := 0
	for _, val := range values {
		sum += val
	}
	return sum
}

func getFirstRepeatedFrequency(values []int) int {
	markers := make(map[int]bool)
	sum := 0
	for i := 0; ; i = (i + 1) % len(values) {
		sum += values[i]
		if _, ok := markers[sum]; ok {
			return sum
		}
		markers[sum] = true
	}
}

func main() {
	values := getInputValues("input.txt")
	fmt.Printf("Resulting frequency = %d\n", getResultingFrequency(values))
	fmt.Printf("First repeated frequency = %d\n", getFirstRepeatedFrequency(values))
}
