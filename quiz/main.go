package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var (
		file_name = flag.String("csv", "problems.csv", "a csv file name")
		limit     = flag.Int("limit", 3, "the time limit for the quiz in secounds")
	)
	flag.Parse()

	problems, score := parseRows(readCsv(file_name)), 0
	ch := make(chan string)
	for i, problem := range problems {
		fmt.Printf("Question #%v: %v\n", i+1, problem.q)
		go userInput(ch)
		select {
		case input := <-ch:
			if input == problem.a {
				score += 1
			}
		case <-time.After(time.Duration(*limit) * time.Second):
			fmt.Println("Time out.")
		}
	}
	fmt.Printf("You scored %v out of %v\n", score, len(problems))
}

type problem struct {
	q string
	a string
}

func parseRows(rows [][]string) []problem {
	restuls := make([]problem, len(rows))
	for i, row := range rows {
		restuls[i] = problem{q: row[0], a: row[1]}
	}
	return restuls
}

func readCsv(file_name *string) [][]string {
	file, err := os.Open(*file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

func userInput(ch chan string) {
	var input string
	fmt.Scanf("%s\n", &input)
	ch <- input
}
