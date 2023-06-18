package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	fileName := flag.String("csv", "problems.csv", "the program accepts a csv file in the format 'question,answer'")
	limit := flag.Int("limit", 30, "set a time limit for the game")
	flag.Parse()

	file, err := os.Open(*fileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file %s", *fileName))
		os.Exit(1)
	}

	reader := csv.NewReader(file)
	questions, err := reader.ReadAll()
	if err != nil {
		exit("Failed to parse CSV file")
	}

	problems := parseLines(questions)
	correct := 0

	timer := time.NewTimer(time.Duration(*limit) * time.Second)

gameLoop:
	for i, problem := range problems {
		fmt.Printf("Question #%d: %s = ", i+1, problem.question)

		answerChannel := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			break gameLoop
		case answer := <-answerChannel:
			if answer == problem.answer {
				correct++
			}
		}
	}

	bound := (correct / len(problems)) * 100
	if correct == len(problems) {
		fmt.Printf("You scored %d out of %d. Way to go. 🥳🥳🥳 \n", correct, len(problems))
	} else if bound >= 80 {
		fmt.Printf("\nYou scored %d out of %d. Good job. 👍\n", correct, len(problems))
	} else {
		fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
	}

}

func parseLines(lines [][]string) []problem {
	parsedlines := make([]problem, len(lines))

	for i, line := range lines {
		parsedlines[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}

	return parsedlines
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
