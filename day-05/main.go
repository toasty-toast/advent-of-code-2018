package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	"unicode"
)

type unit struct {
	Char    rune
	IsUpper bool
}

type node struct {
	Next  *node
	Prev  *node
	Value unit
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadPolymer(filename string) *node {
	file, err := os.Open(filename)
	check(err)

	fscanner := bufio.NewScanner(file)
	fscanner.Scan()
	polymerString := fscanner.Text()

	head := new(node)
	prev := head
	var cur *node
	for _, char := range polymerString {
		var unit unit
		unit.Char = unicode.ToLower(rune(char))
		unit.IsUpper = (unicode.ToUpper(rune(char)) == char)

		cur = new(node)
		cur.Value = unit
		cur.Prev = prev
		prev.Next = cur
		prev = cur
	}

	file.Close()

	return head.Next
}

func copyPolymer(head *node) *node {
	newHead := new(node)
	prev := newHead
	var cur *node
	walkExisting := head
	for walkExisting != nil {
		cur = new(node)
		cur.Value.Char = walkExisting.Value.Char
		cur.Value.IsUpper = walkExisting.Value.IsUpper

		if prev != newHead {
			cur.Prev = prev
		}
		prev.Next = cur
		prev = cur

		walkExisting = walkExisting.Next
	}
	return newHead.Next
}

func reduce(head *node) *node {
	if head == nil {
		return nil
	}

	cur := head
	for cur.Next != nil {
		if (cur.Value.Char == cur.Next.Value.Char) && (cur.Value.IsUpper != cur.Next.Value.IsUpper) {
			if cur == head {
				head = cur.Next.Next
				cur = head
			} else {
				cur.Prev.Next = cur.Next.Next
				if cur.Next.Next != nil {
					cur.Next.Next.Prev = cur.Prev
				}
				cur = cur.Prev
			}
		} else {
			cur = cur.Next
		}
	}

	return head
}

func unitToString(unit unit) string {
	var char rune
	if unit.IsUpper {
		char = unicode.ToUpper(unit.Char)
	} else {
		char = unit.Char
	}
	return string(char)
}

func listToString(head *node) string {
	var builder strings.Builder
	walk := head
	for walk != nil {
		builder.WriteString(unitToString(walk.Value))
		walk = walk.Next
	}
	return builder.String()
}

func removeUnit(head *node, unit rune) *node {
	newHead := head
	cur := head
	for cur != nil {
		if cur.Value.Char == unit {
			if cur == newHead {
				cur.Next.Prev = nil
				newHead = cur.Next
			} else {
				cur.Prev.Next = cur.Next
				if cur.Next != nil {
					cur.Next.Prev = cur.Prev
				}
			}
		}
		cur = cur.Next
	}
	return newHead
}

func removeWorstUnitAndReduce(head *node) *node {
	var bestPolymer *node
	bestLength := math.MaxInt32
	for char := 'a'; char <= 'z'; char++ {
		polymer := copyPolymer(head)
		reduced := reduce(removeUnit(polymer, char))
		reducedLength := len(listToString(reduced))
		if reducedLength < bestLength {
			bestLength = reducedLength
			bestPolymer = reduced
		}
	}
	return bestPolymer
}

func main() {
	head := loadPolymer("input.txt")

	first := reduce(copyPolymer(head))
	fmt.Printf("Length of reduced polymer = %d\n", len(listToString(first)))

	second := removeWorstUnitAndReduce(copyPolymer(head))
	fmt.Printf("Length of reduced polymer after removing worst unit = %d\n", len(listToString(second)))
}
