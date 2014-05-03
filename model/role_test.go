package model

import (
	. "github.com/easykoo/go-blog/common"

	"fmt"
	"testing"
)

type RoleTest struct {
	Id   int    `form:"id" xorm:"int(3) pk not null autoincr"`
	Desc string `form:"description" xorm:"varchar(20) not null"`
}

type UserTest struct {
	Id       int    `form:"id" xorm:"int(11) pk not null autoincr"`
	Username string `form:"username" xorm:"varchar(20) not null"`
	Role     Role   `xorm:"role_id int(3) default 3"`
}

func Test_RoleTest(t *testing.T) {
	SetEngine()
	err := orm.DropTables(&RoleTest{}, &UserTest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = orm.CreateTables(&RoleTest{}, &UserTest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = orm.Insert(&RoleTest{Id: 1, Desc: "test1"}, &UserTest{Id: 1, Username: "username"})
	if err != nil {
		fmt.Println(err)
		return
	}

	userTest := UserTest{}
	_, err = orm.Id(1).Get(&UserTest{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(userTest)

	Expect(t, userTest.Role.Id, 1)
}
