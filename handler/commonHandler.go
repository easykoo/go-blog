package handler

import (
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-blog/common"
	"github.com/easykoo/go-blog/middleware"
	"github.com/easykoo/go-blog/model"
)

func Index(ctx *middleware.Context) {
	ctx.HTML(200, "index", ctx)
}

func About(ctx *middleware.Context) {
	ctx.HTML(200, "about", ctx)
}

func ContactHandler(ctx *middleware.Context, feedback model.Feedback) {
	switch ctx.R.Method {
	case "POST":
		err := feedback.Insert()
		PanicIf(err)
		ctx.Set("success", true)
		ctx.Set("message", Translate(ctx.S.Get("Lang").(string), "message.send.success"))
		ctx.JSON(200, ctx.Response)
	default:
		ctx.HTML(200, "contact", ctx)
	}
}

func LangHandler(ctx *middleware.Context, params martini.Params) {
	lang := params["lang"]
	ctx.S.Set("Lang", lang)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}
