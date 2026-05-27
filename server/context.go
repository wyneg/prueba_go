package server

import (
	"context"
	"encoding/json"
	"net/http"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Cxt            context.Context
	userID         uint
}

func (c *Context) SendText(text string) {
	c.ResponseWriter.Write([]byte(text))
}

func (c *Context) Status(code int) {
	c.ResponseWriter.WriteHeader(code)
}

func (c *Context) JSON(code int, data interface{}) error {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.ResponseWriter.WriteHeader(code)
	return json.NewEncoder(c.ResponseWriter).Encode(data)
}

func (c *Context) BindJSON(dest interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(dest)
}

func (c *Context) SetUserID(id uint) {
	c.userID = id
}

func (c *Context) GetUserID() uint {
	return c.userID
}

func (c *Context) Context() context.Context {
	return c.Cxt
}
