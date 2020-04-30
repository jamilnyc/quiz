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
	"time"
)

func main() {
	pFilename := flag.String("file", "problems.csv", "The name of the CSV file that contains the questions and answers")
	pTimeout := flag.Int("time", 30, "Amount of time to finish the quiz")
	flag.Parse()

	seconds := (*pTimeout)
	ch := make(chan int)
	defer close(ch)

	fmt.Println("Quiz will be from", *pFilename)

	csvFile, err := os.Open(*pFilename)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	scanner := bufio.NewScanner(os.Stdin)
	points := 0
	questionsAnswered := 0

	fmt.Println(seconds, "seconds on the clock")
	fmt.Print("Press <Enter> when you are ready")
	scanner.Scan()
	fmt.Println("Begin!")
	fmt.Println()

	go func() {
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

			questionsAnswered++
			if userInput == answer {
				points++
			}
		}

		ch <- points
	}()

	select {
	case <-ch:
		fmt.Println("You got", points, "out of", questionsAnswered)
	case <-time.After(time.Duration(seconds) * time.Second):
		fmt.Println()
		fmt.Println("You ran out of time!")
		fmt.Println("You got", points, "out of", questionsAnswered)
	}

}
