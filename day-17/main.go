package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

const (
	sand = iota
	clay
	waterSource
	flowingWater
	settledWater
)

func loadData(filename string) [][]int {
	file, _ := os.Open(filename)

	x, y := make([]int, 0), make([]int, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		yRangeRe := regexp.MustCompile(`x=(\d+),\sy=(\d+)\.\.(\d+)`)
		xRangeRe := regexp.MustCompile(`y=(\d+),\sx=(\d+)\.\.(\d+)`)
		if yRangeRe.MatchString(line) {
			match := yRangeRe.FindStringSubmatch(line)
			xVal, _ := strconv.Atoi(match[1])
			yLower, _ := strconv.Atoi(match[2])
			yUpper, _ := strconv.Atoi(match[3])
			for i := yLower; i <= yUpper; i++ {
				x = append(x, xVal)
				y = append(y, i)
			}
		} else if xRangeRe.MatchString(line) {
			match := xRangeRe.FindStringSubmatch(line)
			yVal, _ := strconv.Atoi(match[1])
			xLower, _ := strconv.Atoi(match[2])
			xUpper, _ := strconv.Atoi(match[3])
			y = append(y, yVal)
			y = append(y, yVal)
			x = append(x, xLower)
			x = append(x, xUpper)
			for i := xLower; i <= xUpper; i++ {
				y = append(y, yVal)
				x = append(x, i)
			}
		}
	}

	minX := math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32
	for i := range x {
		if x[i] < minX {
			minX = x[i]
		}
		if x[i] > maxX {
			maxX = x[i]
		}
		if y[i] > maxY {
			maxY = y[i]
		}
	}

	data := make([][]int, 0)
	for i := 0; i <= maxY; i++ {
		data = append(data, make([]int, maxX-minX+3))
	}

	for i := range x {
		data[y[i]][maxX-x[i]+1] = clay
	}

	data[0][maxX-500+1] = waterSource

	file.Close()
	return data
}

func copyData(data [][]int) [][]int {
	duplicate := make([][]int, len(data))
	for i := range data {
		duplicate[i] = make([]int, len(data[i]))
		copy(duplicate[i], data[i])
	}
	return duplicate
}

func spreadWater(data [][]int, x int, y int) {
	maxY := len(data) - 1
	leftClayIndex, rightClayIndex := -1, -1
	for i := x - 1; i >= 0; i-- {
		if data[y][i] == clay {
			leftClayIndex = i
			break
		} else {
			data[y][i] = flowingWater
			if y < maxY && data[y+1][i] == sand {
				break
			}
		}
	}
	for i := x + 1; i < len(data[0]); i++ {
		if data[y][i] == clay {
			rightClayIndex = i
			break
		} else {
			data[y][i] = flowingWater
			if y < maxY && data[y+1][i] == sand {
				break
			}
		}
	}
	if leftClayIndex != -1 && rightClayIndex != -1 {
		for i := leftClayIndex + 1; i < rightClayIndex; i++ {
			data[y][i] = settledWater
		}
	}
}

func doesLowerPieceSupportWater(data [][]int, x int, y int) bool {
	return y != len(data)-1 && (data[y+1][x] == clay || data[y+1][x] == settledWater)
}

func canSpreadLeft(data [][]int, x int, y int) bool {
	return x > 0 && data[y][x-1] == sand
}

func canSpreadRight(data [][]int, x int, y int) bool {
	return x < len(data[0])-1 && data[y][x+1] == sand
}

func shouldSpreadWater(data [][]int, x int, y int) bool {
	return doesLowerPieceSupportWater(data, x, y) && (canSpreadLeft(data, x, y) || canSpreadRight(data, x, y))
	//maxY, maxX := len(data)-1, len(data[0])-1
	// maxX := len(data[0]) - 1
	// if !doesLowerPieceSupportWater(data, x, y) {
	// 	return false
	// }
	// if x > 0 && data[y][x-1] == flowingWater {
	// 	return false
	// }
	// if x < maxX && data[y][x+1] == flowingWater {
	// 	return false
	// }
	// return true
	// return (x > 0 && data[y][x-1] != flowingWater) || doesLowerPieceSupportWater(data, y, x-1)) &&
	// 	(x < maxX && data[y][x+1] != flowingWater) || doesLowerPieceSupportWater(data, y, x+1))
}

func step(data [][]int) [][]int {
	workingData := copyData(data)
	for i := 0; i < len(data)-1; i++ {
		for j := range data[i] {
			if data[i][j] == waterSource || data[i][j] == flowingWater {
				if shouldSpreadWater(data, j, i) {
					spreadWater(workingData, j, i)
				} else if data[i+1][j] == sand {
					workingData[i+1][j] = flowingWater
				}
			}
			if workingData[i][j] == flowingWater && j > 0 && workingData[i][j-1] == clay && j < len(data[i])-1 && workingData[i][j+1] == clay {
				workingData[i][j] = settledWater
			}
		}
	}
	return workingData
}

func printData(data [][]int) {
	for i := range data {
		for j := range data[i] {
			switch data[i][j] {
			case sand:
				fmt.Print(".")
				break
			case clay:
				fmt.Print("#")
				break
			case waterSource:
				fmt.Print("+")
				break
			case flowingWater:
				fmt.Print("|")
				break
			case settledWater:
				fmt.Print("~")
				break
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func countWater(data [][]int) int {
	count := 0
	for i := range data {
		for j := range data[i] {
			if data[i][j] == flowingWater || data[i][j] == settledWater {
				count++
			}
		}
	}
	return count
}

func compareData(data1 [][]int, data2 [][]int) bool {
	for i := range data1 {
		for j := range data1[i] {
			if data1[i][j] != data2[i][j] {
				return false
			}
		}
	}
	return true
}

func main() {
	data := loadData("input.txt")
	for {
		newData := step(data)
		if compareData(data, newData) {
			data = newData
			break
		}
		data = newData
	}

	fmt.Printf("%d tiles have water\n", countWater(data))
}
