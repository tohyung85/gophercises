package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tohyung85/gophercises/exercise-1-quiz-game/quiz"
	"log"
	"os"
	"time"
)

func main() {
	durationPtr := flag.Int("dur", 30, "quiz duration")
	filePath := flag.String("path", "./data/problems-test.csv", "path to csv file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	quizGen := quiz.NewGenerator(file)

	completed := make(chan bool)

	fmt.Println("Press enter to start")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	go quizGen.StartQuiz(completed)
	go startQuizTime(*durationPtr, completed)

	<-completed
	fmt.Printf("Your Score: %d/%d\n", quizGen.CurrentScore, quizGen.TotalQuestions)
}

func startQuizTime(duration int, completed chan bool) {
	fmt.Printf("You have %d seconds \n", duration)
	time.Sleep(time.Duration(duration) * time.Second)
	fmt.Println("Quiz Over - Ran out of time!")
	completed <- true
}
