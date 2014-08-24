package model

import (
	. "github.com/easykoo/go-blog/common"

	"time"
)

type Settings struct {
	Id          int       `form:"settingsId" xorm:"int(11) pk not null autoincr"`
	AppName     string    `form:"appName" xorm:"varchar(45) not null"`
	About       string    `form:"about" xorm:"blob not null"`
	Owner       User      `form:"owner_id" json:"owner_id" xorm:"owner_id"`
	Keywords    string    `form:"keywords" xorm:"varchar(100) default null"`
	Description string    `form:"description" xorm:"varchar(100) default null"`
	CreateUser  string    `xorm:"varchar(20) default null"`
	CreateDate  time.Time `xorm:"datetime created"`
	UpdateUser  string    `xorm:"varchar(20) default null"`
	UpdateDate  time.Time `xorm:"datetime updated"`
	Version     int       `form:"version" xorm:"int(11) version"`
	Page        `xorm:"-"`
}

func GetSettings() *Settings {
	settings := &Settings{Id: 1}
	_, err := orm.Get(settings)
	PanicIf(err)
	return settings
}

func (self *Settings) Insert() error {
	_, err := orm.InsertOne(self)
	Log.Info("Settings ", self.Id, " inserted")
	return err
}

func (self *Settings) Update() error {
	_, err := orm.Id(self.Id).Update(self)
	Log.Info("Settings ", self.Id, " updated!")
	return err
}

func (self *Settings) Delete() error {
	_, err := orm.Delete(self)
	Log.Info("Settings ", self.Id, " deleted")
	return err
}
