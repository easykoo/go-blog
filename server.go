package main

import (
	"github.com/easykoo/binding"
	"github.com/easykoo/sessions"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/russross/blackfriday"

	. "github.com/easykoo/go-blog/auth"
	. "github.com/easykoo/go-blog/common"
	"github.com/easykoo/go-blog/handler"
	"github.com/easykoo/go-blog/middleware"
	"github.com/easykoo/go-blog/model"

	"encoding/gob"
	"html/template"
	"os"
	"time"
)

func init() {
	SetConfig()
	SetLog()
	gob.Register(model.User{})
	gob.Register(model.Settings{})
	Log.Debug("server initializing...")
}

func newMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(middleware.GetLogger())
	m.Map(model.SetEngine())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)

	m.Use(sessions.Sessions("my_session", middleware.NewDbStore(7*24*60*60)))

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Extensions: []string{".tmpl", ".html"},
		Charset:    "UTF-8",
		Funcs: []template.FuncMap{
			{
				"formatTime": func(args ...interface{}) string {
					return args[0].(time.Time).Format("Jan _2 15:04")
				},
				"cnFormatTime": func(args ...interface{}) string {
					return args[0].(time.Time).Format("2006-01-02 15:04")
				},
				"mdToHtml": func(args ...interface{}) template.HTML {
					return template.HTML(string(blackfriday.MarkdownBasic([]byte(args[0].(string)))))
				},
				"unescaped": func(args ...interface{}) template.HTML {
					return template.HTML(args[0].(string))
				},
				"equal": func(args ...interface{}) bool {
					return args[0] == args[1]
				},
				"tsl": func(lang string, format string) string {
					return Translate(lang, format)
				},
				"tslf": func(lang string, format string, args ...interface{}) string {
					return Translatef(lang, format, args...)
				},
				"privilege": func(user interface{}, module int) bool {
					if user == nil {
						return false
					}
					return CheckPermission(user, module)
				},
				"plus": func(args ...int) int {
					var result int
					for _, val := range args {
						result += val
					}
					return result
				},
			},
		},
	}))

	m.Use(middleware.InitContext())
	m.Use(middleware.RecordVisit())

	return &martini.ClassicMartini{m, r}
}

func main() {
	m := newMartini()

	m.Get("/", handler.Blog)
	m.Get("/index", handler.Blog)
	m.Get("/about", handler.About)
	m.Any("/contact", binding.Form(model.Feedback{}), handler.ContactHandler)
	m.Get("/language/change/:lang", handler.LangHandler)

	m.Group("/user", func(r martini.Router) {
		r.Any("/all", AuthRequest(Module_Account), handler.AllUserHandler)
		r.Any("/logout", handler.LogoutHandler)
		r.Any("/login", binding.Form(model.UserLoginForm{}), handler.LoginHandler)
		r.Any("/register", binding.Form(model.UserRegisterForm{}), handler.RegisterHandler)
		r.Any("/delete", AuthRequest(Module_Account), handler.DeleteUsers)
		r.Any("/delete/:id", AuthRequest(Module_Account), handler.DeleteUser)
		r.Any("/role", AuthRequest(Module_Account), handler.SetRole)
		r.Any("/ban/:id", AuthRequest(Module_Account), handler.BanUser)
		r.Any("/lift/:id", AuthRequest(Module_Account), handler.LiftUser)
	})

	m.Group("/profile", func(r martini.Router) {
		r.Any("/profile", AuthRequest(SignInRequired), binding.Form(model.User{}), handler.ProfileHandler)
		r.Any("/preferences", AuthRequest(SignInRequired), handler.PreferencesHandler)
		r.Any("/password", AuthRequest(SignInRequired), binding.Form(model.Password{}), handler.PasswordHandler)
		r.Any("/checkEmail", AuthRequest(SignInRequired), binding.Form(model.User{}), handler.CheckEmail)
	})

	m.Group("/admin", func(r martini.Router) {
		r.Get("/dashboard", AuthRequest(SignInRequired), handler.DashboardHandler)
		r.Any("/settings", AuthRequest(Module_Admin), binding.Form(model.Settings{}), handler.SettingsHandler)
		r.Post("/about", AuthRequest(Module_Admin), handler.AboutHandler)
	})

	m.Group("/feedback", func(r martini.Router) {
		r.Any("/all", AuthRequest(Module_Feedback), handler.AllFeedback)
		r.Any("/info", AuthRequest(Module_Feedback), handler.FeedbackInfo)
		r.Any("/delete", AuthRequest(Module_Feedback), handler.DeleteFeedbackArray)
		r.Any("/delete/:id", AuthRequest(Module_Feedback), handler.DeleteFeedback)
		r.Any("/view/:id", AuthRequest(Module_Feedback), handler.ViewFeedback)
	})

	m.Group("/link", func(r martini.Router) {
		r.Any("/all", AuthRequest(Module_Link), handler.AllLink)
		r.Any("/insert", AuthRequest(Module_Link), binding.Form(model.Link{}), handler.InsertLink)
		r.Any("/delete", AuthRequest(Module_Link), handler.DeleteLinkArray)
		r.Any("/delete/:id", AuthRequest(Module_Link), handler.DeleteLink)
		r.Any("/edit/:id", AuthRequest(Module_Link), handler.EditLink)
	})

	m.Group("/blog", func(r martini.Router) {
		r.Any("", handler.Blog)
		r.Any("/tag/:tag", handler.BlogWithTag)
		r.Any("/view/:id", handler.ViewBlog)
		r.Any("/all", AuthRequest(Module_Blog), handler.AllBlog)
		r.Any("/publish", AuthRequest(Module_Blog), binding.Form(model.Blog{}), handler.PublishBlog)
		r.Any("/save", AuthRequest(Module_Blog), binding.Form(model.Blog{}), handler.SaveBlog)
		r.Any("/edit/:id", AuthRequest(Module_Blog), handler.EditBlog)
		r.Any("/delete", AuthRequest(Module_Blog), handler.DeleteBlogArray)
		r.Any("/delete/:id", AuthRequest(Module_Blog), handler.DeleteBlog)
	})

	m.Group("/blog/comment", func(r martini.Router) {
		r.Any("", handler.Comment)
		r.Any("/delete/:blogId/:seq", AuthRequest(Module_Blog), handler.DeleteComment)
	})

	Log.Info("server is started...")
	os.Setenv("PORT", Cfg.MustValue("", "http_port", "3000"))
	m.Run()
}
