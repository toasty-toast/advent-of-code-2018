package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadValues(filename string) []string {
	file, err := os.Open("input.txt")
	check(err)
	
	values := make([]string, 0)
	fscanner := bufio.NewScanner(file)
    for fscanner.Scan() {
		values = append(values, fscanner.Text())
	}

	return values
}

func checksum(values []string) int {
	twos, threes := 0, 0
	for _, s := range values {
		counters := make(map[rune]int)
		for _, char := range s {
			if _, ok := counters[char]; ok {
				counters[char]++
			} else {
				counters[char] = 1
			}
		}

		gotTwo := false
		gotThree := false
		for _, count := range counters {
			if !gotTwo && count == 2 {
				twos++
				gotTwo = true
			}
			if !gotThree && count == 3 {
				threes++
				gotThree = true
			}
			if gotTwo && gotThree {
				break
			}
		}
	}
	return twos * threes
}

func differByOneChar(first string, second string) bool {
	differences := 0
	for i := 0; i < len(first) && i < len(second); i++ {
		if first[i] != second[i] {
			differences++
		}
		if differences > 1 {
			return false
		}
	}
	return differences == 1
}

func findCorrectIds(values []string) (string, string) {
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if differByOneChar(values[i], values[j]) {
				return values[i], values[j]
			}
		}
	}
	first, second := "", ""
	return first, second
}

func getCommonChars(first string, second string) string {
	var builder strings.Builder
	for i := 0; i < len(first) && i < len(second); i++ {
		if first[i] == second[i] {
			builder.WriteByte(first[i])
		}
	}
	return builder.String()
}

func main() {
	values := loadValues("input.txt")
	fmt.Printf("Checksum = %d\n", checksum(values))
	first, second := findCorrectIds(values)
	fmt.Printf("Correct IDs are %s and %s, common letters are %s\n", first, second, getCommonChars(first, second))
}
