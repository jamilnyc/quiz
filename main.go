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

type problem struct {
	question string
	answer   string
}

func parseQuestionFile(filename string) []problem {
	problems := []problem{}
	fmt.Println("Opening quiz file ", filename, "...")
	csvFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		record, err := reader.Read()

		// Break out of the loop once we've reached the end of the file
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		// Lines should only be a question and answer pair
		if len(record) != 2 {
			continue
		}

		// Clean up the answer to give the user a better chance at
		// correctly answering it
		answer := strings.ToLower(record[1])
		answer = strings.TrimSpace(answer)
		question := record[0]

		p := problem{
			question: question,
			answer:   answer,
		}
		problems = append(problems, p)
	}

	if len(problems) == 0 {
		fmt.Println("There are no valid questions in", filename)
		os.Exit(1)
	}

	return problems
}

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

	scanner := bufio.NewScanner(os.Stdin)
	points := 0
	questionsAnswered := 0
	problems := parseQuestionFile(*pFilename)

	// Randomize the questions if desired
	if *pShuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	fmt.Println(seconds, "seconds on the clock.")
	fmt.Println("There are", len(problems), "questions in total.")
	fmt.Print("Press <Enter> when you are ready ... ")
	scanner.Scan()
	fmt.Println("Begin!")
	fmt.Println()

	go func() {
		for _, p := range problems {
			// prompt user for input
			fmt.Print(p.question, " = ")
			scanner.Scan()
			userInput := scanner.Text()

			// clean up user input
			userInput = strings.ToLower(userInput)
			userInput = strings.TrimSpace(userInput)

			questionsAnswered++
			if userInput == p.answer {
				points++
			}
		}

		// Put any data in the channel to indicate this routine is done
		ch <- points
	}()

	// Run the question/answer routine and a timer concurrently
	// The first one to complete execution will put data in the channel
	select {
	case <-ch:
		fmt.Println()
		fmt.Println("You finished!")
	case <-time.After(time.Duration(seconds) * time.Second):
		fmt.Println()
		fmt.Println()
		fmt.Println("You ran out of time!")
	}

	fmt.Println("Questions Answered:", questionsAnswered, "of", len(problems))
	fmt.Println("Correct Answers:", points)
	percent := 100 * float32(points) / float32(len(problems))
	fmt.Printf("Score: %.1f %%", percent)
}
