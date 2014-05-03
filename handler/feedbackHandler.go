package handler

import (
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-blog/common"
	"github.com/easykoo/go-blog/middleware"
	"github.com/easykoo/go-blog/model"

	"encoding/json"
)

func FeedbackInfo(ctx *middleware.Context) {
	feedback := new(model.Feedback)
	feedback.Viewed = false
	feedbackArray, count, err := feedback.Info()
	PanicIf(err)
	ctx.Set("Array", feedbackArray)
	ctx.Set("Count", count)
	ctx.JSON(200, ctx.Response)
}

func AllFeedback(ctx *middleware.Context) {
	switch ctx.R.Method {
	case "POST":
		feedback := new(model.Feedback)
		feedback.SetPageActive(true)
		feedback.SetPageSize(ParseInt(ctx.R.FormValue("iDisplayLength")))
		feedback.SetDisplayStart(ParseInt(ctx.R.FormValue("iDisplayStart")))
		columnNum := ctx.R.FormValue("iSortCol_0")
		sortColumn := ctx.R.FormValue("mDataProp_" + columnNum)
		feedback.AddSortProperty(sortColumn, ctx.R.FormValue("sSortDir_0"))
		feedbackArray, total, err := feedback.SearchByPage()
		PanicIf(err)
		ctx.Set("aaData", feedbackArray)
		ctx.Set("iTotalDisplayRecords", total)
		ctx.Set("iTotalRecords", total)
		ctx.JSON(200, ctx.Response)
	default:
		ctx.HTML(200, "feedback/allFeedback", ctx)
	}
}

func DeleteFeedback(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	feedback := new(model.Feedback)
	feedback.Id = ParseInt(id)
	err := feedback.Delete()
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func DeleteFeedbackArray(ctx *middleware.Context) {
	feedbackArray := ctx.R.FormValue("feedbackArray")
	feedback := new(model.Feedback)
	var res []int
	json.Unmarshal([]byte(feedbackArray), &res)
	err := feedback.DeleteFeedbackArray(res)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func ViewFeedback(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	feedback := new(model.Feedback)
	feedback.Id = ParseInt(id)
	err := feedback.SetViewed(true)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}
