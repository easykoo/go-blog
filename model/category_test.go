package model

import (
	. "github.com/easykoo/go-blog/common"

	"fmt"
	"testing"
)

func Test_Category(t *testing.T) {
	SetConfig()
	SetLog()
	SetEngine()
	err := orm.DropTables(&Category{})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = orm.CreateTables(&Category{})
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = orm.Insert(&Category{Id: 1, Description: "test1"}, &Category{Id: 2, Description: "test2"}, &Category{Id: 3, Description: "test3"})
	if err != nil {
		fmt.Println(err)
		return
	}

	category := Category{}
	_, err = orm.Id(1).Get(&category)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(category)

	Expect(t, category.Id, 1)
}

func Test_SearchCategory(t *testing.T) {
	SetConfig()
	SetLog()
	SetEngine()
	category := new(Category)
	blogList, total, err := category.SearchByPage()
	Log.Debug(blogList, total, err)
}

func Test_GenerateCategoryId(t *testing.T) {
	SetConfig()
	SetLog()
	SetEngine()
	category := new(Category)
	id, err := category.GenerateCategoryId(0)
	PanicIf(err)

	Expect(t, id, 101)
}
