package model

import (
	. "github.com/easykoo/go-blog/common"

	"time"
)

type SessionInfo struct {
	Id         string    `xorm:"varchar(60) pk not null "`
	Content    string    `xorm:"varchar(15) null"`
	Age        int       `xorm:"int(9)"`
	CreateDate time.Time `xorm:"datetime created"`
	UpdateDate time.Time `xorm:"datetime updated"`
}

func (self *SessionInfo) Exist() (exist bool) {
	exist, err := orm.Get(&SessionInfo{Id: self.Id})
	PanicIf(err)
	return
}

func (self *SessionInfo) GetSessionInfo() (sessionInfo SessionInfo) {
	_, err := orm.Id(self.Id).Get(&sessionInfo)
	PanicIf(err)
	return
}

func (self *SessionInfo) Insert() error {
	_, err := orm.InsertOne(self)
	Log.Infol("SessionInfo ", self.Id, " inserted")
	return err
}

func (self *SessionInfo) Update() error {
	_, err := orm.Id(self.Id).Update(self)
	Log.Infol("SessionInfo ", self.Id, " updated")
	return err
}

func (self *SessionInfo) Delete() error {
	_, err := orm.Delete(self)
	Log.Infol("SessionInfo ", self.Id, " deleted")
	return err
}

func (self *SessionInfo) RemoveExpiredSession() (err error) {
	_, err = orm.Exec("delete from session_info where UNIX_TIMESTAMP(now()) >= age + UNIX_TIMESTAMP(update_date)")
	return
}
