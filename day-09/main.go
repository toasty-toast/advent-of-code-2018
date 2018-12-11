package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type marble struct {
	Value int
	Next  *marble
	Prev  *marble
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadData(filename string) (int, int) {
	file, err := os.Open(filename)
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	file.Close()

	players, _ := strconv.Atoi(strings.Split(line, " ")[0])
	lastMarble, _ := strconv.Atoi(strings.Split(line, " ")[6])
	return players, lastMarble
}

func winningScore(players int, lastMarble int) int {
	scores := make([]int, players)
	cur := new(marble)
	cur.Value = 0
	cur.Next = cur
	cur.Prev = cur
	prev := cur
	for i, marbleScore := 0, 1; marbleScore <= lastMarble; i, marbleScore = (i+1)%players, marbleScore+1 {
		if marbleScore%23 != 0 {
			cur = new(marble)
			cur.Next = prev.Next.Next
			prev.Next.Next.Prev = cur
			cur.Prev = prev.Next
			prev.Next.Next = cur
			cur.Value = marbleScore

			prev = cur
		} else {
			remove := cur.Prev.Prev.Prev.Prev.Prev.Prev.Prev

			scores[i] += marbleScore + remove.Value

			remove.Next.Prev = remove.Prev
			remove.Prev.Next = remove.Next

			prev = remove.Next
		}
	}

	max := 0
	for i := range scores {
		if max < scores[i] {
			max = scores[i]
		}
	}
	return max
}

func main() {
	players, lastMarble := loadData("input.txt")
	fmt.Printf("Winning score: %d\n", winningScore(players, lastMarble))
	fmt.Printf("Winning score in a 100x longer game: %d\n", winningScore(players, lastMarble*100))
}
