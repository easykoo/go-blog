package model

import (
	. "github.com/easykoo/go-blog/common"

	"time"
)

type Visit struct {
	SessionId  string    `xorm:"varchar(60) pk not null "`
	Ip         string    `xorm:"varchar(15) null"`
	User       User      `json:"user_id" xorm:"user_id"`
	CreateDate time.Time `xorm:"datetime created"`
	Page       `xorm:"-"`
}

type Statistics struct {
	TotalCount int
	MonthCount int
	DayCount   int
}

func (self *Visit) Exist() (exist bool) {
	exist, err := orm.Get(self)
	PanicIf(err)
	return
}

func (self *Visit) ExistVisit() (exist bool) {
	exist, err := orm.Get(&Visit{SessionId: self.SessionId})
	PanicIf(err)
	return
}

func (self *Visit) GetVisit() (visit Visit) {
	_, err := orm.Id(self.SessionId).Get(&visit)
	PanicIf(err)
	return
}

func (self *Visit) Insert() error {
	_, err := orm.InsertOne(self)
	Log.Info("Visit ", self.SessionId, " inserted")
	return err
}

func (self *Visit) Update() error {
	_, err := orm.Id(self.SessionId).Update(self)
	Log.Info("Visit ", self.SessionId, " updated")
	return err
}

func (self *Visit) Delete() error {
	_, err := orm.Delete(self)
	Log.Info("Visit ", self.SessionId, " deleted")
	return err
}

func (self *Visit) SearchByPage() ([]Visit, int, error) {
	total, err := orm.Count(self)
	var visit []Visit
	err = orm.OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&visit, self)
	return visit, int(total), err
}

func (self *Visit) GetStatistics() (stat Statistics) {

	total, err := orm.Count(&Visit{})
	PanicIf(err)
	stat.TotalCount = int(total)

	monthCount, err := orm.Where("date_format(create_date,'%Y-%m')=date_format(now(),'%Y-%m')").Count(&Visit{})
	PanicIf(err)
	stat.MonthCount = int(monthCount)

	dayCount, err := orm.Where(" datediff(create_date,NOW()) = 0").Count(&Visit{})
	PanicIf(err)
	stat.DayCount = int(dayCount)

	return
}
