package model

import (
	"time"
)

type Role struct {
	Id          int       `form:"roleId" xorm:"int(3) pk not null autoincr"`
	Description string    `form:"description" xorm:"varchar(20) not null"`
	CreateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	CreateDate  time.Time `xorm:"datetime created"`
	UpdateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	UpdateDate  time.Time `xorm:"datetime updated"`
	Version     int       `form:"version" xorm:"int(11) version"`
}

func (self *Role) GetRoleById(id int) (*Role, error) {
	role := &Role{Id: id}
	_, err := orm.Get(role)
	return role, err
}

type Dept struct {
	Id          int       `form:"deptId" xorm:"int(3) pk not null autoincr"`
	Description string    `form:"description" xorm:"varchar(20) not null"`
	CreateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	CreateDate  time.Time `xorm:"datetime created"`
	UpdateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	UpdateDate  time.Time `xorm:"datetime updated"`
	Version     int       `form:"version" xorm:"int(11) version"`
}

func (self *Dept) GetRoleById(id int) (*Dept, error) {
	dept := &Dept{Id: id}
	_, err := orm.Get(dept)
	return dept, err
}
