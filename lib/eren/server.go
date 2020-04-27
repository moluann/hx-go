package eren

import (
	"net/http"
	"sync"
	"sync/atomic"
)

type Handler interface {
	ServeHTTP(*Context)
}

type HandlerFunc func(*Context)

type HandlersChain []HandlerFunc

//return main handler
func (c HandlersChain) Last() HandlerFunc{
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}


type Engine struct {
	RouterGroup
	pool sync.Pool
	server atomic.Value
	trees methodTrees

	allNoRoute  HandlersChain
	allNoMethod HandlersChain
	noRoute     HandlersChain
	noMethod    HandlersChain
}

//init context
func (engine *Engine) allocContext() *Context{
	return &Context{engine:engine}
}

func (engine *Engine) addRoute(method string, path string, handlers []HandlerFunc) {
	if method == "" {
		panic("eren: HTTP method can not be empty!")
	}
	if path[0] != '/' {
		panic("eren: path must begin with '/'!")
	}
	if len(handlers) == 0 {
		panic("eren: there must be at least one handler!")
	}

	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		engine.trees = append(engine.trees,methodTree{method:method,root:root})
	}
	root.addRoute(path,handlers)
}

func (engine *Engine) Routes()(routes RoutesInfo){
	for _,tree := range engine.trees{
		routes = iterate("",tree.method,routes,tree.root)
	}
	return
}

func iterate(path, method string, routes RoutesInfo, root *node) RoutesInfo {
	path += root.path
	if len(root.handlers) > 0 {
		routes = append(routes, RouteInfo{
			Method:  method,
			Path:    path,
			Handler: nameOfFunction(root.handlers.Last()),
		})
	}
	for _, child := range root.children {
		routes = iterate(path, method, routes, child)
	}
	return routes
}

func NewServer() *Engine{
	engine := &Engine{
		RouterGroup:RouterGroup{
			basePath:	"/",
			Handlers:	nil,
			root:		true,
		},
		trees: 				   make(methodTrees,0,9),
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocContext()
	}
	return engine
}

func (engine *Engine) Start(){

}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter,r *http.Request){

}

