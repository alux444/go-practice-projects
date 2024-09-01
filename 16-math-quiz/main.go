package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseProblem(lines [][]string) []problem {
	//go over lines and parse them into problem struct
	res := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		res[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return res
}

func pullProblems(filePath string) ([]problem, error) {
	//read all problems from questions file
	fileObj, err := os.Open(filePath)
	if err == nil {
		csvReader := csv.NewReader(fileObj)

		csvLines, err := csvReader.ReadAll()
		if err == nil {
			return parseProblem(csvLines), nil
		}

		return nil, fmt.Errorf("Error in reading data in CSV from file: %s - %s", filePath, err.Error())
	}

	return nil, fmt.Errorf("Failed opening file: %s - %s", filePath, err.Error())
}

func main() {
	//input file name
	fname := flag.String("f", "questions.csv", "csv file path")
	//set timer duration
	timer := flag.Int("t", 30, "timer for quiz")
	//pull problems from file
	flag.Parse()
	//handle errors
	problems, err := pullProblems(*fname)
	if err != nil {
		exit(fmt.Sprintf("Something went wrong pulling problems: %s", err.Error()))
	}
	//create variable for counting correct answers
	correctAnswers := 0
	//use timer duration to keep time
	timerObject := time.NewTimer(time.Duration(*timer) * time.Second)
	ansChannel := make(chan string)
	//loop through all problems
problemLoop:
	for i, p := range problems {
		var ans string
		fmt.Printf("Question %d: %s\n", i+1, p.q)

		go func() {
			fmt.Scanf("%s", &ans)
			ansChannel <- ans
		}()
		select {
		case <-timerObject.C:
			fmt.Println("You ran out of time!")
			break problemLoop
		case iAns := <-ansChannel:
			if iAns == p.a {
				correctAnswers++
			}
			if i == len(problems)-1 {
				close(ansChannel)
			}
		}
	}

	fmt.Printf("Score: %d / %d\n", correctAnswers, len(problems))
	<-ansChannel
}
