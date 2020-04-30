package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	pFilename := flag.String("file", "problems.csv", "The name of the CSV file that contains the questions and answers")

	flag.Parse()
	fmt.Println("You want to view", *pFilename)

	csvFile, err := os.Open(*pFilename)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	scanner := bufio.NewScanner(os.Stdin)
	points := 0
	questionCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		question := record[0]
		answer := strings.ToLower(record[1])
		answer = strings.TrimSpace(answer)

		fmt.Print(question, " = ")
		scanner.Scan()
		userInput := scanner.Text()

		// clean up user input
		userInput = strings.ToLower(userInput)
		userInput = strings.TrimSpace(userInput)

		questionCount++
		if userInput == answer {
			points++
		}
	}

	fmt.Println("You got", points, "out of", questionCount)
}
