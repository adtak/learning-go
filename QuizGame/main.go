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
	)
	flag.Parse()

	rows := readCsv(file_name)
	amount, score := len(rows), 0
	ch := make(chan string)

	for i, v := range rows {
		question, answer := v[0], v[1]
		fmt.Printf("Question #%v: %v\n", i+1, question)

		go userInput(ch)
		select {
		case input := <-ch:
			if input == answer {
				score += 1
			}
		case <-time.After(3 * time.Second):
			fmt.Println("Time out.")
		}
	}
	fmt.Printf("You scored %v out of %v\n", score, amount)
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
	fmt.Scan(&input)
	ch <- input
}
