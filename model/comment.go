package model

import (
	. "github.com/easykoo/go-blog/common"

	"time"
)

type Comment struct {
	Blog       Blog      `xorm:"blog_id int(11) pk not null"`
	Seq        int       `xorm:"int(5) pk not null"`
	Name       string    `form:"name" xorm:"varchar(20) null"`
	Www        string    `form:"www" xorm:"varchar(45) null"`
	Email      string    `form:"email" xorm:"varchar(45) null"`
	Content    string    `form:"content" xorm:"varchar(150) not null"`
	ParentSeq  int       `form:"prentSeq" xorm:"int(5) null"`
	Ip         string    `xorm:"varchar(15) null"`
	CreateUser string    `xorm:"varchar(20) default 'SYSTEM'"`
	CreateDate time.Time `xorm:"datetime created"`
	UpdateUser string    `xorm:"varchar(20) default 'SYSTEM'"`
	UpdateDate time.Time `xorm:"datetime updated"`
	Version    int       `form:"version" xorm:"int(11) version"`
	Page       `xorm:"-"`
}

func (self *Comment) GenerateSeq() (int, error) {
	result, err := orm.Query("select max(seq)+1 as seq from comment where blog_id = ?", self.Blog.Id)
	seq := ParseInt(string(result[0]["seq"]))
	if seq < 1 {
		seq = 1
	}
	return seq, err
}

func (self *Comment) Insert() error {
	seq, err := self.GenerateSeq()
	self.Seq = seq
	_, err = orm.InsertOne(self)
	Log.Infol("Comment ", self.Blog.Id, " ", self.Seq, " inserted")
	return err
}

func (self *Comment) Update() error {
	_, err := orm.Update(self)
	Log.Infol("Comment ", self.Blog.Id, " ", self.Seq, " updated")
	return err
}

func (self *Comment) Delete() error {
	_, err := orm.Delete(self)
	Log.Infol("Comment ", self.Blog.Id, " ", self.Seq, " deleted")
	return err
}

func (self *Comment) SearchByPage() ([]Comment, int64, error) {
	total, err := orm.Count(self)
	var comment []Comment
	err = orm.OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&comment, self)
	return comment, total, err
}
