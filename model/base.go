package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	. "github.com/easykoo/go-blog/common"
	"time"
)

var orm *xorm.Engine

func SetEngine() *xorm.Engine {
	Log.Info("db initializing...")
	var err error
	server := Cfg.MustValue("db", "server", "127.0.0.1")
	username := Cfg.MustValue("db", "username", "root")
	password := Cfg.MustValue("db", "password", "pass")
	dbName := Cfg.MustValue("db", "db_name", "go_display")
	orm, err = xorm.NewEngine("mysql", username+":"+password+"@tcp("+server+":3306)/"+dbName+"?charset=utf8")
	PanicIf(err)
	orm.TZLocation = time.Local
	orm.ShowSQL = Cfg.MustBool("db", "show_sql", false)
	orm.Logger = xorm.NewSimpleLogger(Log.GetWriter())
	return orm
}

type DbUtil struct{}

func (self *DbUtil) GetRecentComments() (comments []Comment) {
	err := orm.OrderBy("create_date desc").Limit(5, 0).Find(&comments, &Comment{})
	PanicIf(err)
	return
}

func (self *DbUtil) GetHotBlog() (blog []Blog) {
	result, err := orm.Query("select * from  blog b, (select  blog_id, count(*) count from comment group by blog_id  order by count desc limit 0,5) t where b.id = t.blog_id order by t.count desc, b.create_date desc")
	PanicIf(err)
	for _, val := range result {
		b := Blog{Id: int(ParseInt(string(val["id"]))), Title: string(val["title"])}
		blog = append(blog, b)
	}
	return
}

func (self *DbUtil) GetAllLinks() (links []Link) {
	err := orm.OrderBy("create_date desc").Find(&links, &Link{})
	PanicIf(err)
	return links
}
