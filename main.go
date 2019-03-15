package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ankurgopher/mongobrain/database"
	"github.com/ankurgopher/mongobrain/quiz"
	"github.com/ankurgopher/mongobrain/quizAdmin"
	"github.com/ankurgopher/mongobrain/userManagement"
	"github.com/ankurgopher/mongobrain/util"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"log"
	"net/http"
	"os"
)

func main() {

	commandLine()

	// Goroutine to database connection
	ch := make(chan bool)
	var coll *mongo.Collection

	go func() {
		coll = database.GetQuizCollection()
		ch <- true
	}()

	// Welcome Note
	welcomeNote()

	// Wait for feedback from database
	<-ch

	fmt.Println("You are connected to mongobrain database. You'll need to login.")

	var username string
	var password string
	fmt.Println("Enter your username: ")
	fmt.Scanln(&username)
	fmt.Println("Enter the password: ")
	fmt.Scanln(&password)

	database.GetUserCollection().FindOne(context.TODO(),bson.D{})

	fmt.Println("Starting your quiz.")

	// Start the quiz
	quiz.QuickQuiz(coll)
}

func commandLine() {
	q := flag.Bool("quizzer", false, "Post questions to game for other players")
	flag.Parse()
	if *q == true {
		fmt.Println("Go to http://localhost:8080/questionform to submit questions")
		http.HandleFunc("/", userManagement.SignIn)
		http.HandleFunc("/signup", userManagement.SignUp)
		http.HandleFunc("/signout", userManagement.SignOut)
		http.HandleFunc("/questionform", quizAdmin.RenderQuestion)
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}

func welcomeNote() {
	fmt.Fprintf(os.Stdout, "Welcome to mongobrain.\n Press 1 to start, 0 to exit. \n")

	// Get 1 from user to continue the game or 0 to exit
	var input string
	for {
		_, err := fmt.Scanln(&input)
		if err != nil || (input != "0" && input != "1") {
			msg := fmt.Sprintf("Enter 1 to start, 0 to exit: %s\n", err)
			util.SmallError(msg)
		} else {
			break
		}
	}

	// Exit if 0
	if input == "0" {
		fmt.Println("You chose to exit! See you again. Bye")
		os.Exit(0)
	}
}
