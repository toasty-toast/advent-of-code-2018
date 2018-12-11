package main

import (
	"fmt"
	"math"
)

const puzzleInput = 2187
const gridSize = 300

func powerLevel(row, col, serialNumber int) int {
	return ((((((col + 10) * row) + serialNumber) * (col + 10)) / 100) % 10) - 5
}

func createGrid(serialNumber int) [][]int {
	grid := make([][]int, gridSize)
	for i := range grid {
		grid[i] = make([]int, gridSize)
	}

	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			grid[i][j] = powerLevel(i+1, j+1, serialNumber)
		}
	}

	return grid
}

func subGridPower(grid [][]int, row int, col int, size int) int {
	sum := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			sum += grid[i+row][j+col]
		}
	}
	return sum
}

func largestPower3x3(grid [][]int) (int, int) {
	bestX, bestY, bestSum := 0, 0, math.MinInt32
	for i := 0; i < len(grid)-3; i++ {
		for j := 0; j < len(grid[i])-3; j++ {
			sum := subGridPower(grid, i, j, 3)
			if sum > bestSum {
				bestSum = sum
				bestX = j
				bestY = i
			}
		}
	}
	return bestX + 1, bestY + 1
}

func largestPower(grid [][]int) (int, int, int) {
	bestX, bestY, bestSize, bestSum := 0, 0, 0, math.MinInt32
	for size := 1; size <= 300; size++ {
		for i := 0; i < len(grid)-size; i++ {
			for j := 0; j < len(grid[i])-size; j++ {
				sum := subGridPower(grid, i, j, size)
				if sum > bestSum {
					bestSum = sum
					bestX = j
					bestY = i
					bestSize = size
				}
			}
		}
	}
	return bestX + 1, bestY + 1, bestSize
}

func main() {
	grid := createGrid(puzzleInput)
	{
		row, col := largestPower3x3(grid)
		fmt.Printf("Largest 3x3 square: %d, %d\n", row, col)
	}
	{
		row, col, size := largestPower(grid)
		fmt.Printf("Largest square: %d, %d, %d\n", row, col, size)
	}
}
