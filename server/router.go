package server

import "net/http"

type HandlerFunc func(c *Context)

func (app *App) HttpMethods(method, path string, handler func(*Context)) {
	app.mux.HandleFunc(method+" "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			ResponseWriter: w,
			Request:        r,
			Cxt:            r.Context(),
		})
	})

}
