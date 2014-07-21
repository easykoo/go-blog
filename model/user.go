package model

import (
	"github.com/easykoo/binding"

	. "github.com/easykoo/go-blog/common"

	"net/http"
	"regexp"
	"strconv"
	"time"
)

type User struct {
	Id         int       `form:"userId" xorm:"int(11) pk not null autoincr"`
	Username   string    `form:"username" xorm:"varchar(20) not null"`
	Password   string    `form:"password" xorm:"varchar(60) not null"`
	FullName   string    `form:"fullName" xorm:"varchar(20) null"`
	Gender     int       `form:"gender" xorm:"int(1) default 0"`
	Qq         string    `form:"qq" xorm:"varchar(16) default null"`
	Tel        string    `form:"tel" xorm:"varchar(20) null"`
	Postcode   string    `form:"postcode" xorm:"varchar(10) default null"`
	Address    string    `form:"address" xorm:"varchar(80) default null"`
	Email      string    `form:"email" xorm:"varchar(45) unique"`
	Role       Role      `json:"role_id" xorm:"role_id int(3) default 1"`
	Dept       Dept      `json:"dept_id" xorm:"dept_id int(3) default 1"`
	Active     bool      `xorm:"tinyint(1) default 0"`
	Locked     bool      `xorm:"tinyint(1) default 0"`
	FailTime   int       `xorm:"int(1) default 0"`
	CreateUser string    `xorm:"varchar(20) default 'SYSTEM'"`
	CreateDate time.Time `xorm:"datetime created"`
	UpdateUser string    `xorm:"varchar(20) default 'SYSTEM'"`
	UpdateDate time.Time `xorm:"datetime updated"`
	Version    int       `form:"version" xorm:"int(11) version"`
	Page       `xorm:"-"`
}

func (self *User) ShowName() string {
	if self.FullName != "" {
		return self.FullName
	}
	return self.Username
}

func (self *User) Exist() (bool, error) {
	return orm.Get(self)
}

func (self *User) ExistUsername() (bool, error) {
	return orm.Get(&User{Username: self.Username})
}

func (self *User) ExistEmail() (bool, error) {
	return orm.Get(&User{Email: self.Email})
}

func (self *User) GetUser() (*User, error) {
	user := &User{}
	_, err := orm.Id(self.Id).Get(user)
	return user, err
}

func (self *User) GetUserById(id int) (*User, error) {
	user := &User{Id: id}
	_, err := orm.Get(user)
	return user, err
}

func (self *User) Insert() error {
	self.FullName = self.Username
	self.Dept = Dept{Id: 1}
	self.Role = Role{Id: 4}
	self.Active = true
	self.CreateUser = "SYSTEM"
	self.UpdateUser = "SYSTEM"
	_, err := orm.InsertOne(self)
	Log.Infol(self.Username, " inserted")
	return err
}

func (self *User) Update() error {
	self.Role = Role{}
	self.Dept = Dept{}
	_, err := orm.Id(self.Id).MustCols("gender").Update(self)
	Log.Infol("User ", self.Username, " updated")
	return err
}

func (self *User) Delete() error {
	_, err := orm.Delete(self)
	Log.Infol("User ", self.Username, " deleted")
	return err
}

func (self *User) DeleteUsers(array []int) error {
	_, err := orm.In("id", array).Delete(&User{})
	sql := "delete from `user` where id in ("
	for index, val := range array {
		sql += strconv.Itoa(val)
		if index < len(array)-1 {
			sql += ","
		}
	}
	sql += ")"
	_, err = orm.Exec(sql)
	Log.Infol("Users: ", array, " deleted")
	return err
}

func (self *User) SetRole() error {
	var err error
	_, err = orm.Id(self.Id).MustCols("role_id").Update(&User{Role: self.Role, Version: self.Version})
	Log.Infol("User ", self.Username, " roleId set to ", self.Role.Id)
	return err
}

func (self *User) SetLock(lock bool) error {
	var err error
	self, err = self.GetUser()
	_, err = orm.Id(self.Id).UseBool("locked").Update(&User{Locked: lock, Version: self.Version})
	if lock {
		Log.Infol("User ", self.Username, " locked")
	} else {
		Log.Infol("User ", self.Username, " unlocked")
	}
	return err
}

func (self *User) SelectAll() ([]User, error) {
	var users []User
	err := orm.Find(&users)
	return users, err
}

func (self *User) SearchByPage() ([]User, int64, error) {
	total, err := orm.Count(self)
	var users []User
	err = orm.OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&users, self)
	return users, total, err
}

type UserLoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type UserRegisterForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Email    string `form:"email" binding:"required"`
}

type Password struct {
	Id              int    `form:"id" binding:"required"`
	CurrentPassword string `form:"currentPassword" binding:"required"`
	ConfirmPassword string `form:"confirmPassword" binding:"required"`
}

func (user UserRegisterForm) Validate(errors *binding.Errors, r *http.Request) {
	if len(user.Username) < 5 {
		errors.Fields["username"] = "Length of username should be longer than 5."
	}
	if len(user.Password) < 5 {
		errors.Fields["password"] = "Length of password should be longer than 5."
	}
	re := regexp.MustCompile(".+@.+\\..+")
	matched := re.Match([]byte(user.Email))
	if matched == false {
		errors.Fields["email"] = "Please enter a valid email address."
	}
}

func (password Password) Validate(errors *binding.Errors, r *http.Request) {
	if len(password.ConfirmPassword) < 5 {
		errors.Fields["confirmPassword"] = "Length of password should be longer than 5."
	}
}
