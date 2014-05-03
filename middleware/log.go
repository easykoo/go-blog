package middleware

import (
	"github.com/go-martini/martini"

	. "github.com/easykoo/go-blog/common"

	"net/http"
	"time"
)

func GetLogger() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		start := time.Now()
		Log.Debugf("Started %s %s", req.Method, req.URL.Path)

		rw := res.(martini.ResponseWriter)
		c.Next()

		Log.Debugf("Completed %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}
