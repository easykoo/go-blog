package handler

import (
	"github.com/easykoo/binding"
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-blog/common"
	"github.com/easykoo/go-blog/middleware"
	"github.com/easykoo/go-blog/model"

	"encoding/json"
)

func LogoutHandler(ctx *middleware.Context) {
	ctx.S.Set("SignedUser", nil)
	ctx.Redirect("/index")
}

func LoginHandler(ctx *middleware.Context, formErr binding.Errors, loginUser model.UserLoginForm) {
	switch ctx.R.Method {
	case "POST":
		ctx.JoinFormErrors(formErr)
		password := Md5(loginUser.Password)
		user := &model.User{Username: loginUser.Username, Password: password}
		if !ctx.HasError() {
			if has, err := user.Exist(); has {
				PanicIf(err)
				if user.Locked {
					ctx.Set("User", user)
					ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.invalid.username.or.password"))
					ctx.HTML(200, "user/login", ctx)
					return
				}
				ctx.S.Set("SignedUser", user)
				Log.Info(user.Username, " login")
				ctx.Redirect("/admin/dashboard")
			} else {
				ctx.Set("User", user)
				ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.invalid.username.or.password"))
				ctx.HTML(200, "user/login", ctx)
			}
		} else {
			ctx.HTML(200, "user/login", ctx)
		}
	default:
		ctx.HTML(200, "user/login", ctx)
	}
}

func RegisterHandler(ctx *middleware.Context, formErr binding.Errors, user model.UserRegisterForm) {
	switch ctx.R.Method {
	case "POST":
		ctx.JoinFormErrors(formErr)
		if !ctx.HasError() {
			dbUser := model.User{Username: user.Username, Password: user.Password, Email: user.Email}

			if exist, err := dbUser.ExistUsername(); exist {
				PanicIf(err)
				ctx.AddFieldError("username", Translate(ctx.S.Get("Lang").(string), "message.error.already.exists"))
			}

			if exist, err := dbUser.ExistEmail(); exist {
				PanicIf(err)
				ctx.AddFieldError("email", Translate(ctx.S.Get("Lang").(string), "message.error.already.exists"))
			}

			if !ctx.HasError() {
				dbUser.Password = Md5(user.Password)
				err := dbUser.Insert()
				PanicIf(err)
				ctx.AddMessage(Translate(ctx.S.Get("Lang").(string), "message.register.success"))
			} else {
				ctx.Set("User", user)
			}
			ctx.HTML(200, "user/register", ctx)
		} else {
			ctx.Set("User", user)
			ctx.HTML(200, "user/register", ctx)
		}
	default:
		ctx.HTML(200, "user/register", ctx)
	}
}

func ProfileHandler(ctx *middleware.Context, formErr binding.Errors, user model.User) {
	switch ctx.R.Method {
	case "POST":
		ctx.JoinFormErrors(formErr)
		if !ctx.HasError() {
			err := user.Update()
			PanicIf(err)
			dbUser, err := user.GetUserById(user.Id)
			PanicIf(err)
			ctx.AddMessage(Translate(ctx.S.Get("Lang").(string), "message.change.success"))
			ctx.S.Set("SignedUser", dbUser)
		}
		ctx.HTML(200, "profile/profile", ctx)
	default:
		ctx.HTML(200, "profile/profile", ctx)
	}
}

func PasswordHandler(ctx *middleware.Context, formErr binding.Errors, password model.Password) {
	switch ctx.R.Method {
	case "POST":
		ctx.JoinFormErrors(formErr)
		if !ctx.HasError() {
			if password.CurrentPassword == password.ConfirmPassword {
				ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.password.not.changed"))
			} else {
				user := &model.User{Id: password.Id}
				dbUser, err := user.GetUserById(user.Id)
				PanicIf(err)
				if dbUser.Password == Md5(password.CurrentPassword) {
					dbUser.Password = Md5(password.ConfirmPassword)
					err := dbUser.Update()
					PanicIf(err)
					ctx.AddMessage(Translate(ctx.S.Get("Lang").(string), "message.change.success"))
				} else {
					ctx.AddError(Translate(ctx.S.Get("Lang").(string), "message.error.wrong.password"))
				}
			}
		}
	default:
	}
	ctx.HTML(200, "profile/password", ctx)
}

func CheckEmail(ctx *middleware.Context) {
	if user := ctx.S.Get("SignedUser"); user.(model.User).Email != ctx.R.Form["email"][0] {
		test := &model.User{Email: ctx.R.Form["email"][0]}
		if exist, _ := test.ExistEmail(); exist {
			ctx.JSON(200, Translate(ctx.S.Get("Lang").(string), "message.error.already.exists"))
			return
		}
	}
	ctx.JSON(200, true)
}

func AllUserHandler(ctx *middleware.Context) {
	switch ctx.R.Method {
	case "POST":
		user := new(model.User)
		user.SetPageActive(true)
		user.SetPageSize(ParseInt(ctx.R.FormValue("iDisplayLength")))
		user.SetDisplayStart(ParseInt(ctx.R.FormValue("iDisplayStart")))
		columnNum := ctx.R.FormValue("iSortCol_0")
		sortColumn := ctx.R.FormValue("mDataProp_" + columnNum)
		user.AddSortProperty(sortColumn, ctx.R.FormValue("sSortDir_0"))
		users, total, err := user.SearchByPage()
		PanicIf(err)
		ctx.Set("aaData", users)
		ctx.Set("iTotalDisplayRecords", total)
		ctx.Set("iTotalRecords", total)
		ctx.JSON(200, ctx.Response)
	default:
		ctx.HTML(200, "user/allUser", ctx)
	}
}

func DeleteUser(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	user := new(model.User)
	user.Id = ParseInt(id)
	err := user.Delete()
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func DeleteUsers(ctx *middleware.Context) {
	users := ctx.R.FormValue("Users")
	var res []int
	json.Unmarshal([]byte(users), &res)
	user := new(model.User)
	err := user.DeleteUsers(res)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func SetRole(ctx *middleware.Context) {
	id := ctx.R.PostFormValue("Id")
	roleId := ctx.R.PostFormValue("RoleId")
	version := ctx.R.PostFormValue("Version")
	user := new(model.User)
	user.Id = ParseInt(id)
	user.Role.Id = ParseInt(roleId)
	user.Version = ParseInt(version)
	err := user.SetRole()
	PanicIf(err)
	ctx.Set("success", true)
	Log.Info("User: ", user.Id, " roleId set to ", roleId)
	ctx.JSON(200, ctx.Response)
}

func BanUser(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	user := new(model.User)
	user.Id = ParseInt(id)
	err := user.SetLock(true)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func LiftUser(ctx *middleware.Context, params martini.Params) {
	id := params["id"]
	user := new(model.User)
	user.Id = ParseInt(id)
	err := user.SetLock(false)
	PanicIf(err)
	ctx.Set("success", true)
	ctx.JSON(200, ctx.Response)
}

func PreferencesHandler(ctx *middleware.Context) {
	ctx.HTML(200, "profile/preferences", ctx)
}
