package quiz

import (
	"context"
	"github.com/ankurgopher/mongobrain/database"
	"net/http"
)

func InsertProblem(p Problem) (code int, err error)  {
	coll := database.GetQuizCollection()

	_,err = coll.InsertOne(context.TODO(),p)
	if err!=nil{
		return http.StatusInternalServerError,err
	}
	return http.StatusOK,nil
}