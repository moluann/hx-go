package eren

import (
	"context"
	"math"
	"net/http"
	"net/url"
)

const abortIndex int8 = math.MaxInt8 / 2

func (h HandlerFunc) ServeHTTP(ctx *Context){
	h(ctx)
}

type Context struct {
	context.Context
	Writer http.ResponseWriter
	Request *http.Request
	index int8
	Params Params
	handlers HandlersChain
	Keys map[string]interface{}
	engine *Engine
	errors errors
	Accepted []string
}

func (ctx *Context) reset(){
	ctx.handlers = nil
	ctx.index = -1
	ctx.Keys = nil
	ctx.errors = ctx.errors[0:0]
	ctx.Accepted = nil
}

//return last handlerName
func (ctx *Context) HandlerName() string {
	return nameOfFunction(ctx.handlers.Last())
}

//return main handler
func (ctx *Context) Handler() HandlerFunc {
	return ctx.handlers.Last()
}

/************************************/
/*********** FLOW CONTROL ***********/
/************************************/
func (ctx *Context) Next(){
	ctx.index++
	for s := int8(len(ctx.handlers)); ctx.index < s; ctx.index++ {
		ctx.handlers[ctx.index](ctx)
	}
}

func (ctx *Context) IsAborted() bool{
	return ctx.index >= abortIndex
}

func (ctx *Context) Abort(){
	ctx.index = abortIndex
}

func (ctx *Context) AbortWithStatus(code int) {
	ctx.Status(code)
	ctx.Abort()
}

func (ctx *Context) AbortWithStatusJSON(code int, jsonObj interface{}) {
	ctx.Abort()
	ctx.JSON(code, jsonObj)
}

func (ctx *Context) Status(code int){
	ctx.Writer.WriteHeader(code)
}

func (c *Context) AbortWithError(code int, err error) *Error {
	c.AbortWithStatus(code)
	return c.Error(err)
}

/************************************/
/********* ERROR MANAGEMENT *********/
/************************************/
func (ctx *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	parsedError, ok := err.(*Error)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	ctx.errors = append(ctx.errors, parsedError)
	return parsedError
}

/************************************/
/******** METADATA MANAGEMENT********/
/************************************/

func (ctx *Context) Set(key string, value interface{}) {
	if ctx.Keys == nil {
		ctx.Keys = make(map[string]interface{})
	}
	ctx.Keys[key] = value
}

func (ctx *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = ctx.Keys[key]
	return
}

// GetString returns the value associated with the key as a string.
func (ctx *Context) MustGet(key string) interface{} {
	if value, exists := ctx.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

/************************************/
/************ INPUT DATA ************/
/************************************/
// Param returns the value of the URL param.
// It is a shortcut for c.Params.ByName(key)
//     router.GET("/user/:id", func(c *gin.Context) {
//         // a GET request to /user/john
//         id := c.Param("id") // id == "john"
//     })
func (ctx *Context) Param(key string) string {
	return ctx.Params.ByName(key)
}

// Query returns the keyed url query value if it exists,
// otherwise it returns an empty string `("")`.
// It is shortcut for `c.Request.URL.Query().Get(key)`
//     GET /path?id=1234&name=Manu&value=
// 	   c.Query("id") == "1234"
// 	   c.Query("name") == "Manu"
// 	   c.Query("value") == ""
// 	   c.Query("wtf") == ""
func (ctx *Context) Query(key string) string {
	val,ok := ctx.Request.URL.Query()[key]
	if ok {
		return val[0]
	}
	return ""
}

func (ctx *Context) PostForm(key string){
	ctx.Request.ParseForm()
	ctx.Request.ParseMultipartForm()
}





func (ctx *Context) String(){

}

func (ctx *Context) Json(){

}

func (ctx *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool){
	if path == "" {
		path = "/"
	}
	http.SetCookie(ctx.Writer,&http.Cookie{
		Name:       name,
		Value:      url.QueryEscape(value),
		Path:       path,
		MaxAge:     maxAge,
		Domain:		domain,
		Secure:     secure,
		HttpOnly:   httpOnly,
	})
}

func (ctx *Context) Cookie(name string) (string,error){
	cookie, err := ctx.Request.Cookie(name)
	if err != nil {
		return "",err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeaderNow()
		return
	}

	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}

func (ctx *Context) JSON(code int, obj interface{}) {
	
}








