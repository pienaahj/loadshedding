package main

import (
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

const (
	filename1 string = "data/schedule.csv"
	filename2 string = "data/areas.csv"
	filename3 string = "data/groups.csv"
	filename4 string = "data/schedules.txt"
	filename5 string = "data/schedules.json"
	filename6 string = "data/areas.json"
)

// Declare the start date of the schedule
var day time.Time = time.Date(2020, 8, 20, 0, 0, 0, 0, time.UTC)

//schedule struct
type schedule struct {
	Date  time.Time `json:"date"`
	Stage string    `json:"stage"`
	Group []group   `json:"group"`
}

// group struct
type group struct {
	Group string
	Times []string
}

// area struct
type area struct {
	Group    string   `json:"group"`
	AreaName []string `json:"areaname"`
}

// getArea reads csv file into go struct
func getArea(filename string) []area {
	tmpArea := area{}
	areaX := []area{}
	areas := [][]string{}
	f, fErr := os.Open(filename)
	if fErr != nil {
		log.Fatalf("Cannot open area csv file %v\n", fErr)
	}
	defer f.Close()
	r := csv.NewReader(f)
	areas, csvErr := r.ReadAll()
	if csvErr != nil {
		log.Fatalf("Cannot read csv file %v\n", csvErr)
	}
	fmt.Printf("Completed area read from csv record read%d\n", len(areas))
	//Populate the movies slice
	for _, record := range areas {
		tSlice := []string{}
		tmpArea.Group = record[0]
		for i, v := range record {
			if i != 0 {
				if v != "" {
					tSlice = append(tSlice, v)
				}
			}
		}
		tmpArea.AreaName = tSlice
		areaX = append(areaX, tmpArea)
		tSlice = tSlice[:0]

	}
	return areaX
}

// getGroup reads csv file into go struct
func getGroup(filename string) [][]group {
	tmpGroup := group{}
	groupX := []group{}
	groupXX := make([][]group, 0)
	groups := [][]string{}
	f, fErr := os.Open(filename)
	if fErr != nil {
		log.Fatalf("Cannot open group csv file %v\n", fErr)
	}
	defer f.Close()
	r := csv.NewReader(f)
	groups, csvErr := r.ReadAll()
	if csvErr != nil {
		log.Fatalf("Cannot read csv file %v\n", csvErr)
	}
	fmt.Printf("Completed area read from csv record read: %d\n", len(groups))
	//Populate the group slice
	groupXX = make([][]group, len(groups)) //dimention the stage groups
	for o, record := range groups {        // get the stage
		groupX = groupX[:0]               //clear the group slice
		groupXX[o] = make([]group, 0, 19) //dimention the group slice
		for i, tBlock := range record {   // get the time block
			for _, l := range tBlock { //get the times inside the time block
				if unicode.IsControl(rune(l)) { //discard the control characters
					// fmt.Println("This is a control char: ", l)
					tBlock = strings.TrimFunc(tBlock, func(r rune) bool {
						return !unicode.IsLetter(r) && !unicode.IsNumber(r)
					})
				}
			}
			// Split the time block up into a slice of time seqments & build the group
			tmpGroup.Group = fmt.Sprintf("Group %d", i+1)
			tmpGroup.Times = strings.Split(tBlock, " ")
			// add group to group slice
			groupX = append(groupX, tmpGroup)
			// populate the stage groups
			groupXX[o] = append(groupXX[o], tmpGroup)
		}
		groupX = groupX[:0] //clear the group slice
	}
	return groupXX
}

// writeF writes the output of groups to file
func writeF(schedules []schedule) {
	// Create txt file
	f, err := os.OpenFile(filename4, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("Cannot open txt file! with error: %\n", err)
	}
	defer f.Close()
	fmt.Fprintf(f, "length of schedules: %d\n", len(schedules))
	// Loop through the records
	for i, schedule := range schedules {
		// Print headings
		layout := "Monday, 2 January 2006"
		date := schedule.Date
		t := date.Format(layout)
		fmt.Fprintf(f, "Date: %s	%s\n", t, schedule.Stage)
		// Loop over the 19 groups
		for _, groups := range schedule.Group {
			fmt.Fprintf(f, "record: %d	%s: \nTime:\n", i, groups.Group)
			// Loop over the times
			for _, times := range groups.Times {
				fmt.Fprintf(f, "%s ", times)
				fmt.Fprintln(f)
				fmt.Fprintln(f, strings.Repeat("-", 60))
			}
		}
	}
	fmt.Fprintln(f, strings.Repeat("*", 60))
}

// save the schedule to disk
func dump(filename string, data []schedule) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Cannot create area bin file!: %v\n", err)
	}
	defer f.Close()

	bX, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("Cannot marshal area json: %v\n", err)
	}
	n, err := f.Write(bX)
	if err != nil {
		log.Fatalf("Cannot write area json file!: %v\n", err)
	}
	fmt.Printf("%d lines writen to %s\n", n, filename)
}

// save the schedule to disk
func dumpArea(filename string, data []area) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Cannot create area bin file!: %v\n", err)
	}
	defer f.Close()

	bX, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("Cannot marshal area json: %v\n", err)
	}
	n, err := f.Write(bX)
	if err != nil {
		log.Fatalf("Cannot write area json file!: %v\n", err)
	}
	fmt.Printf("%d lines writen to %s\n", n, filename)
}

// save the schedule to disk
func dumpAreas(filename string, data []area) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Cannot create area bin file!: %v\n", err)
	}
	defer f.Close()
	for _, v := range data {
		err = binary.Write(f, binary.LittleEndian, v)
		if err != nil {
			log.Fatalf("Cannot write area bin file!: %v\n", err)
		}
	}
}

// save the schedule to disk
func dumpSchedules(filename string, data []schedule) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Cannot create schedule bin file!: %v\n", err)
	}
	defer f.Close()
	for _, v := range data {
		err = binary.Write(f, binary.LittleEndian, v)
		if err != nil {
			log.Fatalf("Cannot write schedule bin file!: %v\n", err)
		}
	}
}

// buildSchedule builds the schedule records
func buildSchedule(d time.Time, groups [][]group) []schedule {
	// create a schedule value
	sched := schedule{}
	// create a schedule slice value
	shedules := make([]schedule, 0)
	// set stageCounter
	sCounter := 0
	// Day counter
	dCounter := 0
	// loop through all 19 days for all 6 cycles
	for i := 0; i < 6; i++ { // loop through stages
		// loop through all grougs
		for _, gStage := range groups {
			// check the stage
			stage := ""
			switch sCounter {
			case 0:
				stage = "Stage 1"
			case 1:
				stage = "Stage 2"
			case 2:
				stage = "Stage 3"
			case 3:
				stage = "Stage 4"
			}
			// build the schedule
			sched.Stage = stage

			sched.Group = gStage

			sched.Date = d.Add(time.Duration(dCounter) * time.Hour * 24) //generate the date for stage
			shedules = append(shedules, sched)
			sCounter++ //increment stage counter
			// check stage counter and reset if 4
			if sCounter >= 4 {
				sCounter = 0 //reset stage counter
				dCounter++
			}
		}
	}
	return shedules
}

// test the json file retrieving for schedules
func testJSON(filename string) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Cannot open json file!: %v\n", err)
	}
	defer f.Close()
	scheduleJ := []schedule{}
	bX := []byte{}
	bX, err = ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Cannot read from json file!: %v\n", err)
	}
	err = json.Unmarshal(bX, &scheduleJ)
	if err != nil {
		log.Fatalf("Cannot unmarshal from json file!: %v\n", err)
	}
	for _, v := range scheduleJ {
		fmt.Println(v)
	}
}
func main() {
	// process the areas
	pArea := getArea(filename2)
	dumpArea(filename6, pArea)

	// read the group data
	groupXX := getGroup(filename3)

	// Populate the schedules
	schedules := buildSchedule(day, groupXX)
	writeF(schedules)
	dump(filename5, schedules)

}
