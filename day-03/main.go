package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Rectangle struct {
	Id, Row, Col, Width, Height int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadValues(filename string) []*Rectangle {
	file, err := os.Open("input.txt")
	check(err)
	
	values := make([]*Rectangle, 0)
	fscanner := bufio.NewScanner(file)
    for fscanner.Scan() {
		line := fscanner.Text()
		cols := strings.Split(line, " ")

		rect := new(Rectangle)
		rect.Id, _ = strconv.Atoi(cols[0][1:])
		rect.Col, _ = strconv.Atoi(strings.Split(cols[2], ",")[0])
		rect.Row, _ = strconv.Atoi(strings.Split(strings.Split(cols[2], ",")[1], ":")[0])
		rect.Width, _ = strconv.Atoi(strings.Split(cols[3], "x")[0])
		rect.Height, _ = strconv.Atoi(strings.Split(cols[3], "x")[1])
		
		values = append(values, rect)
	}

	return values
}

func getFabricSize(rects []*Rectangle) (int, int) {
	width, height := 0, 0
	for _, rect := range rects {
		if width < (rect.Col + rect.Width) {
			width = (rect.Col + rect.Width)
		}
		if height < (rect.Row + rect.Height) {
			height = (rect.Row + rect.Height)
		}
	}
	return width, height
}

func createFabric(rects []*Rectangle) [][]int {
	width, height := getFabricSize(rects)
	fabric := make([][]int, height)
	for row := range fabric {
		fabric[row] = make([]int, width)
	}

	for _, rect := range rects {
		for i := 0; i < rect.Height; i++ {
			for j := 0; j < rect.Width; j++ {
				fabric[i + rect.Row][j + rect.Col]++
			}
		}
	}

	return fabric
}

func countContestedSquares(fabric [][]int) int {
	count := 0
	for i := range fabric {
		for j := range fabric[i] {
			if fabric[i][j] > 1 {
				count++
			}
		}
	}

	return count
}

func isRectContested(fabric [][]int, rect *Rectangle) bool {
	for i := 0; i < rect.Height; i++ {
		for j := 0; j < rect.Width; j++ {
			if fabric[i + rect.Row][j + rect.Col] != 1 {
				return true
			}
		}
	}
	return false
}

func findUncontestedRect(fabric [][]int, rects []*Rectangle) *Rectangle {
	for _, rect := range rects {
		if !isRectContested(fabric, rect) {
			return rect
		}
	}
	return nil
}

func main() {
	values := loadValues("input.txt")
	fabric := createFabric(values)
	fmt.Printf("Total contested inches = %d\n", countContestedSquares(fabric))
	fmt.Printf("Uncontested claim ID = %d\n", findUncontestedRect(fabric, values).Id)
}
