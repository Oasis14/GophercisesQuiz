package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// Set the allowed flag params
	csvFilename := flag.String("csv", "problems.csv", "A CSV file formated 'question,answer'")
	timeLimit := flag.Int("limit", 30, "Time limit for quiz in second")
	shuf := flag.Bool("shuf", false, "If true the question order will be shuffled")
	flag.Parse()

	// Open file
	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s", *csvFilename))
	}

	//open file
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse csv file")
	}

	//parse lines to problems struct
	problems := parseLines(lines)

	//Shuffle problems
	shuffle(&problems, *shuf)

	//set up a timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	//loop through each problem for the user
	var correct int
problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			break problemLoop
		case answer := <-answerCh:
			answer = strings.ToLower(answer)
			answer = strings.TrimSpace(answer)
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d of of %d\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.ToLower(strings.TrimSpace(line[1])),
		}
	}
	return ret
}

func shuffle(problems *[]problem, shuf bool) {
	if shuf == true {
		tempProblem := *problems
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(*problems), func(i, j int) { tempProblem[i], tempProblem[j] = tempProblem[j], tempProblem[i] })
		*problems = tempProblem
	}
}

// Create a strct so we arent dependant on what type of file is passed in
type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
