# Quiz Game

## Build
Compile the source as follows:

```
go build -o quiz main.go
```

## Execution

### Sample Run

```
$ ./quiz --time=5
This quiz will be from problems.csv
5 seconds on the clock.
There are 4 questions in total.
Press <Enter> when you are ready ...
Begin!

5+5 = 10
7+3 = 3
1+1 = 2
First U.S. President = 
You ran out of time!
Questions Answered: 3 of 4
Correct Answers: 2
Score: 50.0 %
```

### Help

```
$ ./quiz --help
Usage of C:\path\to\your\executable\quiz:
  -file string
        The name of the CSV file that contains the questions and answers (default "problems.csv")
  -shuffle
        Randomize the questions read from the CSV file
  -time int
        Amount of time to finish the quiz (default 30)
```

## Problems File

The problems file should be a CSV with two columns: question and answer. See below for an example.
Rows that do not contain two columns are ignored.
If there are no valid question/answer pairs, the program terminates.

```
5+5,10
7+3,10
1+1,2
"First U.S. President",George Washington
```