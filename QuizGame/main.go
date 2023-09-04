package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("problems.csv")
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  
  r := csv.NewReader(file)
  rows, err := r.ReadAll()
  if err != nil {
    log.Fatal(err)
  }

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