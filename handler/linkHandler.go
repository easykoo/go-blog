package handler

import (
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-blog/common"
	"github.com/easykoo/go-blog/middleware"
	"github.com/easykoo/go-blog/model"

	"encoding/json"
)

func AllLink(ctx *middleware.Context) {
	switch ctx.R.Method {
	case "POST":
		link := new(model.Link)
		link.SetPageActive(true)
		link.SetPageSize(ParseInt(ctx.R.FormValue("iDisplayLength")))
		link.SetDisplayStart(ParseInt(ctx.R.FormValue("iDisplayStart")))
		columnNum := ctx.R.FormValue("iSortCol_0")
		sortColumn := ctx.R.FormValue("mDataProp_" + columnNum)
		link.AddSortProperty(sortColumn, ctx.R.FormValue("sSortDir_0"))
		linkArray, total, err := link.SearchByPage()
		PanicIf(err)
		ctx.Set("aaData", linkArray)
		ctx.Set("iTotalDisplayRecords", total)
		ctx.Set("iTotalRecords", total)
		ctx.JSON(200, ctx.Response)
	default:
		ctx.HTML(200, "link/allLink", ctx)
	}
}

func InsertLink(ctx *middleware.Context, link model.Link) {
	switch ctx.R.Method {
	case "POST":
		err := link.Insert()
		PanicIf(err)
		ctx.Set("success", true)
		ctx.Set("message", Translate(ctx.S.Get("Lang").(string), "message.send.success"))
		ctx.Redirect("/link/all")
	default:
		ctx.HTML(200, "link/edit", ctx)
	}
}

func DeleteLink(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	link := new(model.Link)
	link.Id = ParseInt(id)
	err := link.Delete()
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func DeleteLinkArray(ctx *middleware.Context) {
	linkArray := ctx.R.FormValue("linkArray")
	link := new(model.Link)
	var res []int
	json.Unmarshal([]byte(linkArray), &res)
	err := link.DeleteLinkArray(res)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func EditLink(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	link := new(model.Link)
	link.Id = ParseInt(id)
	err := link.GetLink()
	PanicIf(err)
	ctx.Set("Link", link)

	ctx.HTML(200, "link/edit", ctx)
}
