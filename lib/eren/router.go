package eren

import "net/http"

type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

type IRoutes interface {
	Use(...HandlerFunc) IRoutes

	Handle(string, string, ...HandlerFunc) IRoutes
	Any(string, ...HandlerFunc) IRoutes
	GET(string, ...HandlerFunc) IRoutes
	POST(string, ...HandlerFunc) IRoutes
	DELETE(string, ...HandlerFunc) IRoutes
	PATCH(string, ...HandlerFunc) IRoutes
	PUT(string, ...HandlerFunc) IRoutes
	OPTIONS(string, ...HandlerFunc) IRoutes
	HEAD(string, ...HandlerFunc) IRoutes

	StaticFile(string, string) IRoutes
	Static(string, string) IRoutes
	StaticFS(string, http.FileSystem) IRoutes
}

type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

func (g *RouterGroup) Use(middleWare ...HandlerFunc) IRoutes {
	g.Handlers = append(g.Handlers,middleWare...)
	return g.returnObj()
}

func (g *RouterGroup) Handle(method string,relativePath string,handlers ...HandlerFunc) IRoutes {
	absolutePath := g.calculateAbsolutePath(relativePath)
	handlers = g.combineHandlers(handlers)
	g.engine.addRoute(method, absolutePath, handlers)
	return g.returnObj()
}

func (g *RouterGroup) Any(string, ...HandlerFunc) IRoutes {
	panic("implement me")
}

func (g *RouterGroup) GET(relativePath string,handlers ...HandlerFunc) IRoutes {
	return g.Handle("GET",relativePath,handlers...)
}

func (g *RouterGroup) POST(relativePath string,handlers ...HandlerFunc) IRoutes {
	return g.Handle("POST",relativePath,handlers...)
}

func (g *RouterGroup) PATCH(relativePath string,handlers ...HandlerFunc) IRoutes {
	return g.Handle("PATCH",relativePath,handlers...)
}

func (g *RouterGroup) DELETE(relativePath string,handlers ...HandlerFunc) IRoutes {
	return g.Handle("DELETE",relativePath,handlers...)
}

func (g *RouterGroup) PUT(relativePath string,handlers ...HandlerFunc) IRoutes {
	return g.Handle("PUT",relativePath,handlers...)
}

func (g *RouterGroup) OPTIONS(relativePath string,handlers ...HandlerFunc) IRoutes {
	return g.Handle("OPTIONS",relativePath,handlers...)
}

func (g *RouterGroup) HEAD(relativePath string,handlers ...HandlerFunc) IRoutes {
	return g.Handle("HEAD",relativePath,handlers...)
}

func (g *RouterGroup) StaticFile(string, string) IRoutes {
	panic("implement me")
}

func (g *RouterGroup) Static(string, string) IRoutes {
	panic("implement me")
}

func (g *RouterGroup) StaticFS(string, http.FileSystem) IRoutes {
	panic("implement me")
}

func (g *RouterGroup) Group(string, ...HandlerFunc) *RouterGroup {
	panic("implement me")
}

var _ IRouter = &RouterGroup{}

type RouteInfo struct {
	Method  string
	Path    string
	Handler string
}

type RoutesInfo []RouteInfo


func (g *RouterGroup) returnObj() IRouter{
	if g.root {
		return g.engine
	}
	return g
}

func (g *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(g.basePath, relativePath)
}

func (g *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(g.Handlers) + len(handlers)
	if finalSize >= int(abortIndex) {
		panic("too many handlers")
	}
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, g.Handlers)
	copy(mergedHandlers[len(g.Handlers):], handlers)
	return mergedHandlers
}