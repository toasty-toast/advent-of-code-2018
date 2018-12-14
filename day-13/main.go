package main

import "fmt"

func main() {
	track := LoadTrack("input.txt")
	ticks := 0
	for {
		//fmt.Printf("%s\n", track)
		colX, colY := track.GetCollision()
		if colX != -1 && colY != -1 {
			fmt.Printf("First collision at %d,%d after %d ticks\n", colX, colY, ticks)
			break
		}
		track.Tick()
		ticks++
	}
}
