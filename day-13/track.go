package main

import (
	"bufio"
	"os"
	"strings"
)

const (
	left = iota
	up
	right
	down
	straight
)

const (
	leftRight = iota
	upDown
	intersection
)

type direction int
type orientation int

type cart struct {
	X, Y                         int
	Direction, NextTurnDirection direction
}

func (c *cart) NextMoveDirection(track *Track) direction {
	optionCount := 0
	var singleOption direction
	if c.X > 0 && track.pieces[c.Y][c.X-1] != nil && track.pieces[c.Y][c.X].Orientation != upDown && track.pieces[c.Y][c.X-1].Orientation != upDown && c.Direction != right {
		optionCount++
		singleOption = left
	}
	if c.X < len(track.pieces[c.Y])-1 && track.pieces[c.Y][c.X+1] != nil && track.pieces[c.Y][c.X].Orientation != upDown && track.pieces[c.Y][c.X+1].Orientation != upDown && c.Direction != left {
		optionCount++
		singleOption = right
	}
	if c.Y > 0 && track.pieces[c.Y-1][c.X] != nil && track.pieces[c.Y][c.X].Orientation != leftRight && track.pieces[c.Y-1][c.X].Orientation != leftRight && c.Direction != down {
		optionCount++
		singleOption = up
	}
	if c.Y < len(track.pieces)-1 && track.pieces[c.Y+1][c.X] != nil && track.pieces[c.Y][c.X].Orientation != leftRight && track.pieces[c.Y+1][c.X].Orientation != leftRight && c.Direction != up {
		optionCount++
		singleOption = down
	}

	if optionCount == 1 {
		c.Direction = singleOption
		return singleOption
	}

	c.Turn()
	return c.Direction
}

func (c *cart) Turn() {
	if c.NextTurnDirection == left {
		switch c.Direction {
		case left:
			c.Direction = down
			break
		case up:
			c.Direction = left
			break
		case right:
			c.Direction = up
			break
		case down:
			c.Direction = right
			break
		}
		c.NextTurnDirection = straight
	} else if c.NextTurnDirection == straight {
		c.NextTurnDirection = right
	} else if c.NextTurnDirection == right {
		switch c.Direction {
		case left:
			c.Direction = up
			break
		case up:
			c.Direction = right
			break
		case right:
			c.Direction = down
			break
		case down:
			c.Direction = left
			break
		}
		c.NextTurnDirection = left
	}
}

func (c *cart) Move(track *Track, direction direction) {
	newX, newY := 0, 0
	switch direction {
	case left:
		newX = c.X - 1
		newY = c.Y
		break
	case up:
		newX = c.X
		newY = c.Y - 1
		break
	case right:
		newX = c.X + 1
		newY = c.Y
		break
	case down:
		newX = c.X
		newY = c.Y + 1
		break
	}

	track.pieces[c.Y][c.X].Cart = nil
	track.pieces[newY][newX].Cart = c
	c.X = newX
	c.Y = newY
}

type trackPiece struct {
	Cart        *cart
	Orientation orientation
}

func (t *trackPiece) String() string {
	if t.Cart != nil {
		switch t.Cart.Direction {
		case left:
			return "<"
		case right:
			return ">"
		case up:
			return "^"
		case down:
			return "v"
		}
	}
	switch t.Orientation {
	case leftRight:
		return "-"
	case upDown:
		return "|"
	case intersection:
		return "+"
	}
	panic("Invalid track piece configuration")
}

// Track represents a track that carts can move on.
type Track struct {
	pieces [][]*trackPiece
	carts  []*cart
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// LoadTrack creates a Track from an input file.
func LoadTrack(filename string) *Track {
	file, err := os.Open(filename)
	check(err)
	scanner := bufio.NewScanner(file)

	track := new(Track)
	track.pieces = make([][]*trackPiece, 0)
	track.carts = make([]*cart, 0)
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]*trackPiece, len(line))
		for i := range line {
			if line[i] != ' ' {
				row[i] = new(trackPiece)
				if line[i] == '<' || line[i] == '>' || line[i] == '^' || line[i] == 'v' {
					row[i].Cart = new(cart)
					row[i].Cart.Y = lineNum
					row[i].Cart.X = i
					row[i].Cart.NextTurnDirection = left
					track.carts = append(track.carts, row[i].Cart)
					switch line[i] {
					case '<':
						row[i].Cart.Direction = left
						row[i].Orientation = leftRight
						break
					case '>':
						row[i].Cart.Direction = right
						row[i].Orientation = leftRight
						break
					case '^':
						row[i].Cart.Direction = up
						row[i].Orientation = upDown
						break
					case 'v':
						row[i].Cart.Direction = down
						row[i].Orientation = upDown
						break
					}
				} else {
					switch line[i] {
					case '|':
						row[i].Orientation = upDown
						break
					case '-':
						row[i].Orientation = leftRight
						break
					case '\\':
						fallthrough
					case '/':
						fallthrough
					case '+':
						row[i].Orientation = intersection
						break
					}
				}
			}
		}
		lineNum++
		track.pieces = append(track.pieces, row)
	}

	return track
}

// Tick advances each cart on the track one space.
func (t *Track) Tick() {
	for i := range t.carts {
		nextMoveDir := t.carts[i].NextMoveDirection(t)
		t.carts[i].Move(t, nextMoveDir)
	}
}

func (t *Track) String() string {
	var builder strings.Builder
	for i := range t.pieces {
		for j := range t.pieces[i] {
			if t.pieces[i][j] != nil {
				builder.WriteString(t.pieces[i][j].String())
			} else {
				builder.WriteString(" ")
			}
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

// GetCollision returns the X, Y coordinates of the first collision.
func (t *Track) GetCollision() (int, int) {
	for i := 0; i < len(t.carts); i++ {
		for j := i + 1; j < len(t.carts); j++ {
			if t.carts[i].X == t.carts[j].X && t.carts[i].Y == t.carts[j].Y {
				return t.carts[i].X, t.carts[i].Y
			}
		}
	}
	return -1, -1
}
