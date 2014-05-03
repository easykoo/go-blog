package model

import (
	. "github.com/easykoo/go-blog/common"

	"testing"
)

func Test_SearchTag(t *testing.T) {
	SetConfig()
	SetLog()
	SetEngine()
	result, err := orm.Query("select name, count(name)  as count from tag group by name order by count desc")
	PanicIf(err)
	var tagInfoArray []TagInfo
	for _, val := range result {
		tagInfoArray = append(tagInfoArray, TagInfo{Name: string(val["name"]), Count: ParseInt(string(val["count"]))})
	}

	Log.Debug(tagInfoArray)
}
