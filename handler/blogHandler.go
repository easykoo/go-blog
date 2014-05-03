package handler

import (
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-blog/common"
	"github.com/easykoo/go-blog/middleware"
	"github.com/easykoo/go-blog/model"

	"encoding/json"
	"time"
)

func PublishBlog(ctx *middleware.Context, blog model.Blog) {
	switch ctx.R.Method {
	case "POST":
		if blog.Title == "" || blog.Content == "" {
			ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.publish.failed"))
		} else {
			tags := ctx.R.PostForm["tags"]
			blog.SetTags(tags)
			signedUser := ctx.S.Get("SignedUser").(model.User)
			blog.State = "PUBLISHED"
			blog.PublishDate = time.Now()
			if blog.Version == 0 {
				blog.Priority = 5
				blog.Author = signedUser
				blog.CreateUser = signedUser.Username
				err := blog.Insert()
				PanicIf(err)
			} else {
				err := blog.Update()
				PanicIf(err)
			}
		}
		ctx.Redirect("/blog/view/" + IntString(blog.Id))
	default:
		tags, err := blog.GetAllTags()
		PanicIf(err)
		ctx.Set("Tags", tags)

		ctx.HTML(200, "blog/edit", ctx)
	}
}

func SaveBlog(ctx *middleware.Context, blog model.Blog) {
	if blog.Title == "" || blog.Content == "" {
		ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.save.failed"))
	} else {
		tags := ctx.R.PostForm["tags"]
		blog.SetTags(tags)
		signedUser := ctx.S.Get("SignedUser").(model.User)
		blog.UpdateUser = signedUser.Username
		if blog.Version == 0 {
			blog.State = "DRAFT"
			blog.Priority = 5
			blog.Author = signedUser
			blog.CreateUser = signedUser.Username
			err := blog.Insert()
			PanicIf(err)
		} else {
			err := blog.Update()
			PanicIf(err)
		}
		dbBlog, err := blog.GetBlogById()
		PanicIf(err)
		ctx.Set("Blog", dbBlog)

		ctx.AddMessage(Translate(ctx.S.Get("Lang").(string), "message.save.success"))
	}

	tags, err := blog.GetAllTags()
	PanicIf(err)
	ctx.Set("Tags", tags)

	ctx.HTML(200, "blog/edit", ctx)
}

func AllBlog(ctx *middleware.Context) {
	switch ctx.R.Method {
	case "POST":
		blog := new(model.Blog)
		blog.SetPageActive(true)
		blog.SetPageSize(ParseInt(ctx.R.FormValue("iDisplayLength")))
		blog.SetDisplayStart(ParseInt(ctx.R.FormValue("iDisplayStart")))
		columnNum := ctx.R.FormValue("iSortCol_0")
		sortColumn := ctx.R.FormValue("mDataProp_" + columnNum)
		blog.AddSortProperty(sortColumn, ctx.R.FormValue("sSortDir_0"))
		blogList, total, err := blog.SearchByPage(false)
		PanicIf(err)
		ctx.Set("aaData", blogList)
		ctx.Set("iTotalDisplayRecords", total)
		ctx.Set("iTotalRecords", total)
		ctx.JSON(200, ctx.Response)
	default:
		ctx.HTML(200, "blog/allBlog", ctx)
	}
}

func Blog(ctx *middleware.Context) {
	blog := new(model.Blog)
	blog.SetPageActive(true)
	blog.SetPageSize(10)
	pageNo := ParseInt(ctx.R.FormValue("page"))
	blog.SetPageNo(pageNo)
	blog.State = "PUBLISHED"
	blog.AddSortProperty("publish_date", "desc")
	blogList, total, err := blog.SearchByPage(true)
	PanicIf(err)

	blog.SetTotalRecord(total)
	blog.Result = blogList
	ctx.Set("Blog", blog)

	tags, err := blog.GetAllTags()
	PanicIf(err)
	ctx.Set("Tags", tags)

	ctx.HTML(200, "blog", ctx)
}

func BlogWithTag(ctx *middleware.Context, params martini.Params) {
	tagName := params["tag"]
	blog := new(model.Blog)
	blog.SetPageActive(true)
	blog.SetPageSize(10)
	pageNo := ParseInt(ctx.R.FormValue("page"))
	blog.SetPageNo(pageNo)
	blog.State = "PUBLISHED"
	blog.AddSortProperty("publish_date", "desc")
	blogList, total, err := blog.SearchWithTagByPage(tagName)
	PanicIf(err)
	ctx.Set("Tag", tagName)

	blog.SetTotalRecord(total)
	blog.Result = blogList
	ctx.Set("Blog", blog)
	tags, err := blog.GetAllTags()
	PanicIf(err)
	ctx.Set("Tags", tags)

	ctx.HTML(200, "blog", ctx)
}

func ViewBlog(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	blog := new(model.Blog)
	blog.Id = ParseInt(id)
	err := blog.GetBlog()
	blog.Visit += 1
	blog.UpdateVisit()
	PanicIf(err)
	ctx.Set("Blog", blog)

	tags, err := blog.GetAllTags()
	PanicIf(err)
	ctx.Set("Tags", tags)

	ctx.HTML(200, "blog/view", ctx)
}

func DeleteBlog(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	blog := new(model.Blog)
	blog.Id = ParseInt(id)
	err := blog.Delete()
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func DeleteBlogArray(ctx *middleware.Context) {
	blogArray := ctx.R.FormValue("blogArray")
	blog := new(model.Blog)
	var res []int
	json.Unmarshal([]byte(blogArray), &res)
	err := blog.DeleteBlogArray(res)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func EditBlog(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	blog := new(model.Blog)
	blog.Id = ParseInt(id)
	err := blog.GetBlog()
	PanicIf(err)
	ctx.Set("Blog", blog)

	tags, err := blog.GetAllTags()
	PanicIf(err)
	ctx.Set("Tags", tags)

	ctx.HTML(200, "blog/edit", ctx)
}

func Comment(ctx *middleware.Context) {
	id := ParseInt(ctx.R.PostFormValue("blogId"))
	name := ctx.R.PostFormValue("name")
	email := ctx.R.PostFormValue("email")
	www := ctx.R.PostFormValue("www")
	content := ctx.R.PostFormValue("content")
	comment := model.Comment{Blog: model.Blog{Id: id}, Name: name, Email: email, Www: www, Content: content}

	if comment.Content == "" {
		ctx.Set("success", false)
		ctx.Set("message", Translate(ctx.S.Get("Lang").(string), "message.error.submit.failed"))
	} else {
		comment.Ip = GetRemoteIp(ctx.R)
		err := comment.Insert()
		PanicIf(err)
		ctx.Set("success", true)
		ctx.Set("message", Translate(ctx.S.Get("Lang").(string), "message.submit.success"))
	}
	ctx.JSON(200, ctx.Response)
}

func DeleteComment(ctx *middleware.Context, params martini.Params) {
	blogId := ParseInt(params["blogId"])
	seq := ParseInt(params["seq"])
	comment := model.Comment{Blog: model.Blog{Id: blogId}, Seq: seq}

	err := comment.Delete()
	if err != nil {
		ctx.Set("success", false)
		ctx.Set("message", Translate(ctx.S.Get("Lang").(string), "message.error.delete.failed"))
	} else {
		ctx.Set("success", true)
		ctx.Set("message", Translate(ctx.S.Get("Lang").(string), "message.delete.success"))
	}
	ctx.JSON(200, ctx.Response)
}
