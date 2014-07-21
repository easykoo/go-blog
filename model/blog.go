package model

import (
	. "github.com/easykoo/go-blog/common"

	"strconv"
	"strings"
	"time"
)

type Blog struct {
	Id          int       `form:"blogId" xorm:"int(11) pk not null autoincr"`
	Title       string    `form:"title" xorm:"varchar(45) not null"`
	Content     string    `form:"content" xorm:"blob not null"`
	State       string    `xorm:"varchar(10) default null"`
	Priority    int       `xorm:"int(1) default 5"`
	Author      User      `json:"author_id" xorm:"author_id"`
	Visit       int       `xorm:"int(9)"`
	Tags        []Tag     `form:"tags" json:"tags" xorm:"-"`
	Comments    []Comment `json:"comments" xorm:"-"`
	PublishDate time.Time `xorm:"datetime default null"`
	CreateUser  string    `xorm:"varchar(20) default null"`
	CreateDate  time.Time `xorm:"datetime created"`
	UpdateUser  string    `xorm:"varchar(20) default null"`
	UpdateDate  time.Time `xorm:"datetime updated"`
	Version     int       `form:"version" xorm:"int(11) version"`
	Page        `xorm:"-"`
}

type Tag struct {
	Name string `xorm:"varchar(20) pk not null"`
	Blog Blog   `xorm:"blog_id int(11) pk not null"`
}

type TagInfo struct {
	Name  string
	Count int
}

func (self *Blog) Insert() error {
	session := orm.NewSession()
	defer session.Close()
	_, err := session.InsertOne(self)
	for key, _ := range self.Tags {
		self.Tags[key].Blog.Id = self.Id
		session.Insert(self.Tags[key])
	}
	Log.Infol("Blog ", self.Id, " inserted")
	return err
}

func (self *Blog) SetTags(tagNames []string) {
	tagNames = strings.Split(tagNames[0], ",")
	var tags []Tag
	for _, val := range tagNames {
		tag := Tag{Name: val, Blog: Blog{Id: self.Id}}
		tags = append(tags, tag)
	}
	self.Tags = tags
}

func (self *Blog) LoadTagsFromDb() {
	tag := &Tag{Blog: Blog{Id: self.Id}}
	var tags []Tag
	err := orm.Omit("blog_id").Find(&tags, tag)
	PanicIf(err)
	self.Tags = tags
}

func (self *Blog) GetComments() []Comment {
	comment := &Comment{Blog: Blog{Id: self.Id}}
	var comments []Comment
	err := orm.Omit("blog_id").Find(&comments, comment)
	PanicIf(err)
	return comments
}

func (self *Blog) GetCommentSize() (size int64) {
	comment := &Comment{Blog: Blog{Id: self.Id}}
	size, err := orm.Count(comment)
	PanicIf(err)
	return
}

func (self *Blog) Update() error {
	session := orm.NewSession()
	defer session.Close()
	_, err := session.Id(self.Id).Update(self)
	_, err = session.Exec("delete from tag where blog_id = ?", self.Id)
	for key, _ := range self.Tags {
		session.Insert(self.Tags[key])
	}
	Log.Infol("Blog ", self.Id, " updated!")
	return err
}

func (self *Blog) UpdateVisit() error {
	_, err := orm.Id(self.Id).Cols("visit").Update(self)
	Log.Infol("Blog ", self.Id, " updated!")
	return err
}

func (self *Blog) Delete() error {
	session := orm.NewSession()
	defer session.Close()
	_, err := session.Delete(self)
	_, err = session.Exec("delete from comment where blog_id = ?", self.Id)
	_, err = session.Exec("delete from tag where blog_id = ?", self.Id)
	Log.Infol("Blog ", self.Id, " deleted")
	return err
}

func (self *Blog) DeleteTags() error {
	_, err := orm.Exec("delete from tag where blog_id = ?", self.Id)
	return err
}

func (self *Blog) GetBlogById() (*Blog, error) {
	blog := &Blog{Id: self.Id}
	_, err := orm.Get(blog)
	return blog, err
}

func (self *Blog) GetBlog() error {
	_, err := orm.Id(self.Id).Get(self)
	return err
}

func (self *Blog) SetState(state string) error {
	var err error
	_, err = orm.Id(self.Id).MustCols("state").Update(&Blog{State: state})
	return err
}

func (self *Blog) DeleteBlogArray(array []int) error {
	_, err := orm.In("id", array).Delete(&Blog{})
	sql := "delete from `blog` where id in ("
	for index, val := range array {
		sql += strconv.Itoa(val)
		if index < len(array)-1 {
			sql += ","
		}
	}
	sql += ")"
	_, err = orm.Exec(sql)
	Log.Infol("Blog Array: ", array, " deleted")
	return err
}

func (self *Blog) SearchByPage(content bool) ([]Blog, int, error) {
	total, err := orm.Count(self)
	var blog []Blog
	if content {
		err = orm.OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&blog, self)
	} else {
		err = orm.Omit("content").OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&blog, self)
	}
	return blog, int(total), err
}

func (self *Blog) SearchWithTagByPage(tag string) ([]Blog, int, error) {
	total, err := orm.Join("LEFT", "tag", "blog.id=tag.blog_id").Where("tag.name=?", tag).Count(self)
	var blog []Blog
	err = orm.Join("LEFT", "tag", "blog.id=tag.blog_id").Where("tag.name=?", tag).OrderBy(self.GetSortProperties()[0].Column+" "+self.GetSortProperties()[0].Direction).Limit(self.GetPageSize(), self.GetDisplayStart()).Find(&blog, self)
	return blog, int(total), err
}

func BatchLoadTagsFromDb(blog []Blog) {
	for key, _ := range blog {
		blog[key].LoadTagsFromDb()
	}
}

func (self *Blog) GetAllTags() ([]TagInfo, error) {
	result, err := orm.Query("select name, count(name)  as count from tag group by name order by count desc")
	PanicIf(err)
	var tagInfoArray []TagInfo
	for _, val := range result {
		tagInfoArray = append(tagInfoArray, TagInfo{Name: string(val["name"]), Count: ParseInt(string(val["count"]))})
	}
	return tagInfoArray, err
}

func (self *Tag) GetBlogByTag() ([]Blog, error) {
	var tags []Tag
	err := orm.Find(&tags, self)
	PanicIf(err)
	var blog []Blog
	for _, val := range tags {
		blog = append(blog, val.Blog)
	}
	return blog, err
}

func (self *Blog) GetTags() []Tag {
	var tags []Tag
	err := orm.Find(&tags, &Tag{Blog: Blog{Id: self.Id}})
	PanicIf(err)
	return tags
}

func (self *Blog) Summary() string {
	return strings.Split(self.Content, "----------")[0]
}

func (self *Blog) AllContent() string {
	return strings.Replace(self.Content, "----------", "", -1)
}
