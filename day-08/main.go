package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type node struct {
	Children []*node
	Metadata []int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadNode(data []int, index int) (*node, int) {
	numChildren := data[index]
	numMetadata := data[index+1]

	consumedPoints := 2
	root := new(node)
	for i := 0; i < numChildren; i++ {
		node, consumed := loadNode(data, index+consumedPoints)
		root.Children = append(root.Children, node)
		consumedPoints += consumed
	}
	for i := 0; i < numMetadata; i++ {
		root.Metadata = append(root.Metadata, data[index+consumedPoints+i])
	}
	return root, consumedPoints + numMetadata
}

func treeFromData(data []int) *node {
	root, _ := loadNode(data, 0)
	return root
}

func loadTree(filename string) *node {
	file, err := os.Open(filename)
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	stringVals := strings.Split(scanner.Text(), " ")
	data := make([]int, len(stringVals))
	for i := range stringVals {
		data[i], _ = strconv.Atoi(stringVals[i])
	}
	root := treeFromData(data)

	file.Close()
	return root
}

func sumMetadata(node *node) int {
	sum := 0
	for i := range node.Metadata {
		sum += node.Metadata[i]
	}
	for i := range node.Children {
		sum += sumMetadata(node.Children[i])
	}
	return sum
}

func nodeValue(node *node) int {
	sum := 0
	if len(node.Children) > 0 {
		for i := range node.Metadata {
			index := node.Metadata[i] - 1
			if index >= 0 && index < len(node.Children) {
				sum += nodeValue(node.Children[index])
			}
		}
	} else {
		sum += sumMetadata(node)
	}
	return sum
}

func main() {
	root := loadTree("input.txt")
	fmt.Printf("Sum of all metadata entries: %d\n", sumMetadata(root))
	fmt.Printf("Value of root node: %d\n", nodeValue(root))
}
