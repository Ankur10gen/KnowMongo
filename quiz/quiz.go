package quiz

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ankurgopher/mongobrain/util"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)


type Problem struct {
	Question string `bson:"q"`
	Choices []string `bson:"c"`
	Answer uint8 `bson:"a"`
	Category uint8 `bson:"category"` // 0 uncategorised 1 basic 2 intermediate 3 advanced
	Author string `bson:"author"`
	SubmittedOn time.Time `bson:"time"`
	Verify uint8 `bson:"verify"` // 0 unverified 1 verified
}

var score int
var totalQues int

func presentProblem(q string, c []string, a uint8) {

	// Print question
	fmt.Println(q)

	// Print choices
	for i,v := range c{
		fmt.Printf("%d: %s \t",i+1,v)
	}
	fmt.Println()
	fmt.Println("Enter your answer:")

	var max int
	max = len(c)

	// Listen for answer
	var answer uint8

	for {
		var s string
		fmt.Scanln(&s)
		ai,err := strconv.Atoi(s)
		if err!=nil || ai > max{
			fmt.Println("Option doesn't exist. Enter again.")
		} else {
			answer = uint8(ai)
			if answer == a+1{
				score++
			}
			break
		}
	}
}

// StartQuiz starts the quiz
// It gets 10 questions from the collection
func QuickQuiz(coll *mongo.Collection)  {
	//get 10 questions randomly
	pipeline := bson.A{
		bson.D{
			{
				"$sample",
				bson.D{
					{
						"size",
						10,
					},
				},
			},
		},
	}
	cur, err := coll.Aggregate(context.TODO(),pipeline)

	if err!=nil{
		msg := fmt.Sprintf("Couldn't fetch quiz %s\n",err)
		util.BigError(msg)
	}

	for cur.Next(context.TODO()){
		var p Problem
		cur.Decode(&p)
		presentProblem(p.Question,p.Choices,p.Answer)
		totalQues++
		fmt.Println()
	}

	fmt.Printf("You scored %d out of %d.\n",score,totalQues)
}