package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Time time.Time
	Text string
}

type GuardInfo struct {
	Id        int
	SleepData [60]int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func timeFromString(timeString string) time.Time {
	year, _ := strconv.Atoi(strings.Split(timeString, "-")[0])
	month, _ := strconv.Atoi(strings.Split(timeString, "-")[1])
	day, _ := strconv.Atoi(strings.Split(strings.Split(timeString, "-")[2], " ")[0])
	hour, _ := strconv.Atoi(strings.Split(strings.Split(timeString, " ")[1], ":")[0])
	minute, _ := strconv.Atoi(strings.Split(strings.Split(timeString, " ")[1], ":")[1])
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
}

func loadRecords(filename string) []*Record {
	file, err := os.Open("input.txt")
	check(err)

	records := make([]*Record, 0)
	fscanner := bufio.NewScanner(file)
	for fscanner.Scan() {
		line := fscanner.Text()

		timeString := strings.Split(line, "]")[0][1:]
		textString := strings.Split(line, "]")[1][1:]

		record := new(Record)
		record.Time = timeFromString(timeString)
		record.Text = textString
		records = append(records, record)
	}

	file.Close()

	sort.Slice(records, func(i, j int) bool {
		return records[i].Time.Before(records[j].Time)
	})

	return records
}

func calculateGuardInfo(records []*Record) []*GuardInfo {
	guards := make([]*GuardInfo, 0)
	dict := make(map[int]*GuardInfo)
	var guard *GuardInfo
	for i := 0; i < len(records); i++ {
		if strings.Contains(records[i].Text, "Guard #") {
			guardId, _ := strconv.Atoi(strings.Split(strings.Split(records[i].Text, "#")[1], " ")[0])
			if _, ok := dict[guardId]; ok {
				guard = dict[guardId]
			} else {
				guard = new(GuardInfo)
				guard.Id = guardId
				dict[guardId] = guard
				guards = append(guards, guard)
			}
		} else if strings.Contains(records[i].Text, "falls asleep") {
			startTime := records[i].Time.Minute()
			endTime := records[i+1].Time.Minute()
			for j := startTime; j < endTime; j++ {
				guard.SleepData[j]++
			}
		}
	}

	return guards
}

func totalMinutesAsleep(guard *GuardInfo) int {
	minutes := 0
	for i := range guard.SleepData {
		minutes += guard.SleepData[i]
	}
	return minutes
}

func mostMinutesAsleep(guards []*GuardInfo) *GuardInfo {
	var guard *GuardInfo
	for i := range guards {
		if guard == nil || totalMinutesAsleep(guards[i]) > totalMinutesAsleep(guard) {
			guard = guards[i]
		}
	}
	return guard
}

func bestMinute(guard *GuardInfo) int {
	minute := 0
	for i := range guard.SleepData {
		if guard.SleepData[minute] < guard.SleepData[i] {
			minute = i
		}
	}
	return minute
}

func bestFrequency(guards []*GuardInfo) (*GuardInfo, int) {
	var guard *GuardInfo
	var minute int
	for i := range guards {
		bestMinute := bestMinute(guards[i])
		if guard == nil || guard.SleepData[minute] < guards[i].SleepData[bestMinute] {
			guard = guards[i]
			minute = bestMinute
		}
	}
	return guard, minute
}

func main() {
	records := loadRecords("input.txt")
	guardInfo := calculateGuardInfo(records)
	{
		guard := mostMinutesAsleep(guardInfo)
		minute := bestMinute(guard)
		fmt.Printf("Most time asleep: ID = %d, Best minute = %d, Solution = %d\n", guard.Id, minute, guard.Id*minute)
	}
	{
		guard, minute := bestFrequency(guardInfo)
		fmt.Printf("Most frequently asleep: ID = %d, Minute = %d, Solution = %d\n", guard.Id, minute, guard.Id*minute)
	}
}
