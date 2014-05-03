package middleware

import (
	"github.com/easykoo/sessions"
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-blog/common"
	"github.com/easykoo/go-blog/model"

	"net/http"
	"reflect"
)

func RecordVisit() martini.Handler {
	return func(s sessions.Session, r *http.Request) {
		visit := new(model.Visit)
		visit.SessionId = s.GetId()
		user := s.Get("SignedUser")
		var id int
		if user != nil {
			if reflect.TypeOf(user).Kind() == reflect.Struct {
				id = user.(model.User).Id
			} else {
				id = user.(*model.User).Id
			}
		}
		visit.User = model.User{Id: id}
		visit.Ip = GetRemoteIp(r)
		if visit.ExistVisit() {
			visit.Update()
		} else {
			visit.Insert()
		}
	}
}
