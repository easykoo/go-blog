package model

import (
	. "github.com/easykoo/go-blog/common"

	"testing"
)

func Init() {
	SetConfig()
	SetLog()
	SetEngine()
}

func Test_GetHotBlog(t *testing.T) {
	Init()
	blog := new(DbUtil).GetHotBlog()
	Log.Debug(blog)
}
