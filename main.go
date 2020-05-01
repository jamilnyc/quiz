package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// Define command line flags
	pFilename := flag.String("file", "problems.csv", "The name of the CSV file that contains the questions and answers")
	pTimeout := flag.Int("time", 30, "Amount of time to finish the quiz")
	pShuffle := flag.Bool("shuffle", false, "Randomize the questions read from the CSV file")
	flag.Parse()

	seconds := (*pTimeout)

	// Channel that will be written to if the user finishes on time
	ch := make(chan int)
	defer close(ch)

	fmt.Println("This quiz will be from", *pFilename)

	csvFile, err := os.Open(*pFilename)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	scanner := bufio.NewScanner(os.Stdin)
	points := 0
	questionsAnswered := 0

	var lines [][]string
	for {
		record, err := reader.Read()

		// Break out of the loop once we've reached the end of the file
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if len(record) != 2 {
			continue
		}

		// Clean up the answer to give the user a better chance at
		// correctly answering it
		answer := strings.ToLower(record[1])
		answer = strings.TrimSpace(answer)
		record[1] = answer

		lines = append(lines, record)
	}
	totalQuestions := len(lines)

	if totalQuestions == 0 {
		fmt.Println("There are no valid questions in", *pFilename)
		os.Exit(1)
	}

	// Randomize the questions if desired
	if *pShuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(lines), func(i, j int) {
			lines[i], lines[j] = lines[j], lines[i]
		})
	}

	fmt.Println(seconds, "seconds on the clock.")
	fmt.Println("There are", totalQuestions, "questions in total.")
	fmt.Print("Press <Enter> when you are ready ... ")
	scanner.Scan()
	fmt.Println("Begin!")
	fmt.Println()

	go func() {
		for _, line := range lines {
			question, answer := line[0], line[1]
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
		fmt.Println()
		fmt.Println("You finished!")
	case <-time.After(time.Duration(seconds) * time.Second):
		fmt.Println()
		fmt.Println("You ran out of time!")
	}

	fmt.Println("Questions Answered:", questionsAnswered, "of", totalQuestions)
	fmt.Println("Correct Answers:", points)
	percent := 100 * float32(points) / float32(totalQuestions)
	fmt.Printf("Score: %.1f %%", percent)
}
