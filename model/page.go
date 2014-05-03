package model

import (
	. "github.com/easykoo/go-blog/common"
)

type SortProperty struct {
	Column    string
	Direction string
}

type Page struct {
	pageActive     bool
	displayStart   int
	pageNo         int
	pageSize       int
	totalRecord    int
	totalPage      int
	sortProperties []*SortProperty
	Result         interface{}
}

func (self *Page) GetDisplayStart() int {
	return self.displayStart
}

func (self *Page) GetPageNo() int {
	return self.pageNo
}

func (self *Page) GetPreviousPageNo() int {
	if self.pageNo <= 1 {
		return 1
	}
	return self.pageNo - 1
}

func (self *Page) GetNextPageNo() int {
	if self.pageNo >= self.totalPage {
		return self.totalPage
	}
	return self.pageNo + 1
}

func (self *Page) GetPageSize() int {
	return self.pageSize
}

func (self *Page) GetTotalRecord() int {
	return self.totalRecord
}

func (self *Page) GetTotalPage() int {
	return self.totalPage
}

func (self *Page) GetSortProperties() []*SortProperty {
	return self.sortProperties
}

func (self *Page) SetPageActive(active bool) {
	self.pageActive = active
}

func (self *Page) SetPageSize(pageSize int) {
	self.pageSize = pageSize
}

func (self *Page) SetDisplayStart(displayStart int) {
	self.initIf()
	if ((displayStart + 1) / self.pageSize) < 1 {
		self.pageNo = 1
	} else {
		self.pageNo = ((displayStart + 1) / self.pageSize) + 1
	}
	self.displayStart = displayStart
}

func (self *Page) SetPageNo(pageNo int) {
	if pageNo < 1 {
		pageNo = 1
	}
	self.pageNo = pageNo
	self.displayStart = (pageNo - 1) * self.pageSize
}

func (self *Page) SetTotalRecord(totalRecord int) {
	self.initIf()
	self.totalRecord = totalRecord
	var totalPage int
	if totalRecord%self.pageSize == 0 {
		totalPage = totalRecord / self.pageSize
	} else {
		totalPage = totalRecord/self.pageSize + 1
	}
	self.totalPage = totalPage
	if self.pageNo > totalPage {
		self.pageNo = totalPage
	}
}

func (self *Page) AddSortProperty(column string, direction string) {
	self.sortProperties = append(self.sortProperties, &SortProperty{Column: Atoa(column), Direction: direction})
}

func (self *Page) initIf() {
	if self.pageNo == 0 {
		self.pageNo = 1
	}
	if self.pageSize == 0 {
		self.pageSize = 10
	}
}
