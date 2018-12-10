package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type point struct {
	PosX, PosY, VelX, VelY int
}

type bounds struct {
	X, Y, Width, Height int
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func (b bounds) Area() int64 {
	return int64(b.Width) * int64(b.Height)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadPoints(filename string) []*point {
	file, err := os.Open(filename)
	check(err)

	points := make([]*point, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		re := regexp.MustCompile("-?[0-9]+")
		values := re.FindAllString(line, -1)

		newPoint := new(point)
		newPoint.PosX, _ = strconv.Atoi(values[0])
		newPoint.PosY, _ = strconv.Atoi(values[1])
		newPoint.VelX, _ = strconv.Atoi(values[2])
		newPoint.VelY, _ = strconv.Atoi(values[3])
		points = append(points, newPoint)
	}

	file.Close()
	return points
}

func clonePoints(points []*point) []*point {
	newPoints := make([]*point, len(points))
	for i, oldPoint := range points {
		newPoint := new(point)
		newPoint.PosX = oldPoint.PosX
		newPoint.PosY = oldPoint.PosY
		newPoint.VelX = oldPoint.VelX
		newPoint.VelY = oldPoint.VelY
		newPoints[i] = newPoint
	}
	return newPoints
}

func boundingBox(points []*point) bounds {
	minX, minY, maxX, maxY := math.MaxInt32, math.MaxInt32, math.MinInt32, math.MinInt32
	for i := range points {
		if minX > points[i].PosX {
			minX = points[i].PosX
		}
		if maxX < points[i].PosX {
			maxX = points[i].PosX
		}
		if minY > points[i].PosY {
			minY = points[i].PosY
		}
		if maxY < points[i].PosY {
			maxY = points[i].PosY
		}
	}
	var boundingBox bounds
	boundingBox.X = minX
	boundingBox.Y = minY
	boundingBox.Width = abs(maxX - minX)
	boundingBox.Height = abs(maxY - minY)
	return boundingBox
}

func timeStepForward(points []*point) {
	for i := range points {
		points[i].PosX += points[i].VelX
		points[i].PosY += points[i].VelY
	}
}

func timeStepBackward(points []*point) {
	for i := range points {
		points[i].PosX -= points[i].VelX
		points[i].PosY -= points[i].VelY
	}
}

func pointOccupied(x, y int, points []*point) bool {
	for i := range points {
		if points[i].PosX == x && points[i].PosY == y {
			return true
		}
	}
	return false
}

func stringFromPoints(points []*point) string {
	boundingBox := boundingBox(points)
	var builder strings.Builder
	for i := boundingBox.Y; i <= boundingBox.Y+boundingBox.Height; i++ {
		for j := boundingBox.X; j <= boundingBox.X+boundingBox.Width; j++ {
			if pointOccupied(j, i, points) {
				builder.WriteString("#")
			} else {
				builder.WriteString(".")
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

func message(points []*point) (string, int) {
	lastArea := int64(math.MaxInt64)
	for i := 0; ; i++ {
		boundingBox := boundingBox(points)
		area := boundingBox.Area()
		if lastArea < area {
			timeStepBackward(points)
			return stringFromPoints(points), i - 1
		}
		lastArea = area
		timeStepForward(points)
	}
}

func main() {
	points := loadPoints("input.txt")
	message, seconds := message(points)
	fmt.Printf("Message:\n%s", message)
	fmt.Printf("Message appeared after %ds\n", seconds)
}
