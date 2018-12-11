package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type coordinate struct {
	ID int
	X  int
	Y  int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadCoordinates(filename string) []*coordinate {
	file, err := os.Open(filename)
	check(err)

	id := 0
	coordinates := make([]*coordinate, 0)
	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		line := fscanner.Text()
		coordinate := new(coordinate)
		id++
		coordinate.ID = id
		coordinate.X, _ = strconv.Atoi(strings.Split(line, ",")[0])
		coordinate.Y, _ = strconv.Atoi(strings.Split(line, " ")[1])
		coordinates = append(coordinates, coordinate)
	}

	file.Close()
	return coordinates
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func manhattanDistance(x int, y int, coordinate *coordinate) int {
	return abs(coordinate.X-x) + abs(coordinate.Y-y)
}

func getClosestCoordinate(x int, y int, coordinates []*coordinate) *coordinate {
	var closest *coordinate
	closestDist := math.MaxInt32
	ties := 0
	for _, coordinate := range coordinates {
		dist := manhattanDistance(x, y, coordinate)
		if dist == closestDist {
			ties++
		}
		if dist < closestDist {
			ties = 0
			closestDist = dist
			closest = coordinate
		}
	}

	if ties > 0 {
		return nil
	}
	return closest
}

func mapSize(coordinates []*coordinate) (int, int) {
	width, height := 0, 0
	for _, coordinate := range coordinates {
		if width < coordinate.X {
			width = coordinate.X
		}
		if height < coordinate.Y {
			height = coordinate.Y
		}
	}
	return width, height
}

func createMapArray(coordinates []*coordinate) [][]int {
	width, height := mapSize(coordinates)

	arr := make([][]int, width+1)
	for i := range arr {
		arr[i] = make([]int, height+1)
	}

	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			closest := getClosestCoordinate(i, j, coordinates)
			if closest == nil {
				arr[i][j] = -1
			} else {
				arr[i][j] = closest.ID
			}
		}
	}

	return arr
}

func isAreaInfinite(id int, arr [][]int) bool {
	height := len(arr)
	width := len(arr[0])
	for i := 0; i < height; i++ {
		if arr[i][0] == id {
			return true
		}
		if arr[i][width-1] == id {
			return true
		}
	}
	for i := 0; i < width; i++ {
		if arr[0][i] == id {
			return true
		}
		if arr[height-1][i] == id {
			return true
		}
	}
	return false
}

func largestFiniteArea(arr [][]int) int {
	counter := make(map[int]int)
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			value := arr[i][j]
			if _, contains := counter[value]; !contains {
				counter[value] = 1
			} else {
				counter[value]++
			}
		}
	}

	largest := 0
	for id, area := range counter {
		if !isAreaInfinite(id, arr) && largest < area {
			largest = area
		}
	}
	return largest
}

func safeRegionSize(coordinates []*coordinate) int {
	width, height := mapSize(coordinates)
	count := 0
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			totalDist := 0
			for _, coordinate := range coordinates {
				totalDist += manhattanDistance(i, j, coordinate)
			}
			if totalDist < 10000 {
				count++
			}
		}
	}
	return count
}

func main() {
	coordinates := loadCoordinates("input.txt")
	arr := createMapArray(coordinates)
	fmt.Printf("Larest finite area = %d\n", largestFiniteArea(arr))
	fmt.Printf("Safe region size = %d\n", safeRegionSize(coordinates))
}
