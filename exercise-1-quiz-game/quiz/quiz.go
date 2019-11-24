package quiz

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

type Quiz struct {
	p string
	a string
}

type QuizGenerator struct {
	csvReader      *csv.Reader
	TotalQuestions int
	CurrentScore   int
}

func NewGenerator(f *os.File) *QuizGenerator {
	return &QuizGenerator{
		csvReader:      csv.NewReader(f),
		TotalQuestions: 0,
		CurrentScore:   0,
	}
}

func (qg *QuizGenerator) StartQuiz(completed chan bool) {
	records, err := qg.csvReader.ReadAll()

	if err != nil {
		log.Fatal(err)
		completed <- true
		return
	}

	qg.TotalQuestions = len(records)
	fmt.Println("And here we go...")

	for _, entry := range records {
		quiz := parseCSVRecord(entry)
		fmt.Println(quiz.p)
		reader := bufio.NewReader(os.Stdin)
		inp, _ := reader.ReadString('\n')
		inp = strings.TrimSpace(inp)
		if inp == strings.TrimSpace(quiz.a) {
			qg.CurrentScore += 1
			fmt.Println("You are correct!")
		} else {
			fmt.Println("No...that's not right...")
		}
		fmt.Printf("Your Answer: %s\n", inp)
		fmt.Printf("Right Answer: %s\n", entry[1])
	}
}

func parseCSVRecord(rec []string) *Quiz {
	quiz := Quiz{p: rec[0], a: rec[1]}
	return &quiz
}
