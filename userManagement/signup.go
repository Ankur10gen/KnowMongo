package userManagement

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/ankurgopher/mongobrain/database"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Name     string `bson:"name"`
	Password string `bson:"password"`
}

type userCreation struct {
	Err bool
	Success bool
}

var t *template.Template
var e userCreation

func init()  {
	t = template.Must(t.ParseFiles("templates/userManagement/signup.html"))
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t.ExecuteTemplate(w, "signup",nil)
	} else if r.Method == http.MethodPost {
		//t, err := template.ParseFiles("templates/userManagement/signup.html")
		//if err != nil {
		//	http.Redirect(w, r, "/", http.StatusInternalServerError)
		//}
		//t.Execute(w, nil)
		r.ParseForm()
		var u User
		for k, v := range r.Form {
			if k == "username" {
				u.Name = strings.Join(v, "")
			} else if k == "password" {
				data := strings.Join(v,"")
				hp := fmt.Sprintf("%x", md5.Sum([]byte(data)))
				u.Password = hp
			}
		}

		coll := database.GetUserCollection()
		_, err := coll.InsertOne(context.TODO(), u)
		if err != nil {
			log.Println(err)
			e.Err = true
			t.ExecuteTemplate(w,"signup", e)
		} else {
			e.Err = false
			e.Success = true
			t.ExecuteTemplate(w,"signup",e)
		}
	}
}
