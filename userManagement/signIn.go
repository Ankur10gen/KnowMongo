package userManagement

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/ankurgopher/mongobrain/database"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/satori/go.uuid"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type userLogin struct {
	Fail bool
}

func init()  {
	t = template.Must(t.ParseFiles("templates/userManagement/signin.html"))
}

func verifyUser(u User) bool {

	filter := bson.D{
		{
			"name",
			u.Name,
		},
		{
			"password",
			u.Password,
		},
	}

	coll := database.GetUserCollection()
	res := coll.FindOne(context.TODO(),filter)

	var ph User
	err := res.Decode(&ph)
	log.Println(err)

	if err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func SignIn(w http.ResponseWriter, r *http.Request)  {
	if r.Method == http.MethodGet{
		t.ExecuteTemplate(w,"signin",nil)
	} else if r.Method == http.MethodPost{
		r.ParseForm()
		var u User
		for k,v := range r.Form{
			if k == "username" {
				u.Name = strings.Join(v, "")
			} else if k == "password" {
				data := strings.Join(v,"")
				hp := fmt.Sprintf("%x", md5.Sum([]byte(data)))
				u.Password = hp
			}
		}

		if verifyUser(u){
			// Set cookie
			sID,_ := uuid.NewV4()
			c := &http.Cookie{
				Name:"session-mb",
				Value:sID.String(),
			}
			http.SetCookie(w,c)

			coll := database.GetSessionsCollection()
			coll.DeleteMany(context.TODO(),bson.D{{"username",u.Name}})
			coll.InsertOne(context.TODO(),bson.D{{"username",u.Name},{"sID",sID.String()}})

			http.Redirect(w,r,"/questionform",http.StatusFound)
			return
		} else {
			var msg userLogin
			msg.Fail = true
			t.ExecuteTemplate(w,"signin",msg)
		}
	}
}
