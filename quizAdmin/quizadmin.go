package quizAdmin

import (
	"context"
	"fmt"
	"github.com/ankurgopher/mongobrain/util"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ankurgopher/mongobrain/database"
	"github.com/ankurgopher/mongobrain/quiz"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
)

type data struct {
	Message string
	Author string
}

var a data

type session struct {
	Username string `bson:"username"`
	SID string `bson:"sID"`
}


func RenderQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c, err := r.Cookie("session-mb")
		if err == http.ErrNoCookie{
			http.Redirect(w,r,"/",http.StatusSeeOther)
			return
		}

		coll := database.GetSessionsCollection()
		filter := bson.D{
			{
				"sID",
				c.Value,
			},
		}
		sr := coll.FindOne(context.TODO(),filter)
		var res session
		err = sr.Decode(&res)

		if err == mongo.ErrNoDocuments {
				http.Redirect(w,r,"/",http.StatusSeeOther)
				return
		}

		t, err := template.ParseFiles("quizAdmin/setQuestion.html")
		if err!=nil{
			util.SmallError(err.Error())
		}

		a.Author = res.Username
		t.ExecuteTemplate(w,"setQuestion", a)

	} else if r.Method == http.MethodPost {
		var prob quiz.Problem
		r.ParseForm()
		for k, v := range r.Form {
			if k == "question" {
				prob.Question = strings.Join(v, "")
			} else if k == "c" {
				prob.Choices = v
			} else if k == "answer" {
				a, _ := strconv.Atoi(strings.Join(v, ""))
				prob.Answer = uint8(a) - 1
			} else if k == "category" {
				c, _ := strconv.Atoi(strings.Join(v, ""))
				prob.Category = uint8(c)
			}
		}
		prob.SubmittedOn = time.Now()
		prob.Author = a.Author

		code, err := quiz.InsertProblem(prob)
		if err != nil {
			http.Error(w, "Couldn't post question.", code)
			return
		}

		if code == http.StatusOK{
			t, err := template.ParseFiles("quizAdmin/setQuestion.html")
			if err!=nil{
				util.SmallError(err.Error())
			}
			msg := fmt.Sprintf("Last submission was successful at %s",time.Now().String())
			a.Message = msg
			t.ExecuteTemplate(w,"setQuestion",a)
		}
	}
}
