package model

import (
	. "github.com/easykoo/go-blog/common"
	"time"
)

type Link struct {
	Id          int       `form:"id" xorm:"int(3) pk not null autoincr"`
	Description string    `form:"description" xorm:"varchar(40) not null"`
	Url         string    `form:"url" xorm:"varchar(80) not null"`
	CreateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	CreateDate  time.Time `xorm:"datetime created"`
	UpdateUser  string    `xorm:"varchar(20) default 'SYSTEM'"`
	UpdateDate  time.Time `xorm:"datetime updated"`
	Version     int       `form:"version" xorm:"int(11) version"`
	Page        `xorm:"-"`
}

func (self *Link) Insert() error {
	session := orm.NewSession()
	defer session.Close()
	_, err := session.InsertOne(self)
	Log.Info("Link ", self.Id, " inserted")
	return err
}

func (self *Link) Update() error {
	session := orm.NewSession()
	defer session.Close()
	_, err := session.Id(self.Id).Update(self)
	Log.Info("Link ", self.Id, " updated!")
	return err
}

func (self *Link) Delete() error {
	session := orm.NewSession()
	defer session.Close()
	_, err := session.Delete(self)
	Log.Info("Link ", self.Id, " deleted")
	return err
}

func (self *Link) GetLinkById() (*Link, error) {
	link := &Link{Id: self.Id}
	_, err := orm.Get(link)
	return link, err
}

func (self *Link) GetLink() error {
	_, err := orm.Id(self.Id).Get(self)
	return err
}

func (self *Link) DeleteLinkArray(array []int) error {
	_, err := orm.In("id", array).Delete(&Link{})
	Log.Info("Link Array: ", array, " deleted")
	return err
}

func (self *Link) SearchByPage() ([]Link, int, error) {
	total, err := orm.Count(self)
	var link []Link
	err = orm.OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&link, self)
	return link, int(total), err
}
