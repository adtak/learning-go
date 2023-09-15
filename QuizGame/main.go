package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var (
		file_name = flag.String("csv", "problems.csv", "a csv file name")
	)
	flag.Parse()
	rows := readCsv(file_name)
	var score, amount int32
	for i, v := range rows {
		question, answer := v[0], v[1]
		fmt.Printf("Question #%v: %v\n", i+1, question)
		var input string
		fmt.Scan(&input)
		if input == answer {
			score += 1
		}
		amount += 1
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