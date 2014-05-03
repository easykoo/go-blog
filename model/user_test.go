package model

import (
	. "github.com/easykoo/go-blog/common"

	"testing"
)

func Test_user(t *testing.T) {
	SetEngine()
	user := &User{Username: "test4", Password: "11111", Email: "ddd3@ddd.com"}
	user.Delete()
	err := user.Insert()
	PanicIf(err)

	dbUser, err1 := user.GetUser()
	PanicIf(err1)
	Expect(t, dbUser.Dept.Id, 1)
	Expect(t, dbUser.Role.Id, 3)
}
