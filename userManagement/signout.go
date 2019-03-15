package userManagement

import "net/http"

func SignOut(w http.ResponseWriter, r *http.Request)  {
	c := &http.Cookie{Name:"session-mb",Value:"reset"}
	http.SetCookie(w,c)
	http.Redirect(w,r,"/",http.StatusSeeOther)
	return
}
