package tiny_gin

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc defines the request handler used by TinyGin
type HandlerFunc func(ctx *Context)

// Engine implement the interface of ServerHTTP
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup

	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

// New is the constructor of Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// Run defines the method to start a http server
func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, r)
	c.engine = e
	c.handlers = middlewares
	e.router.handle(c)
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	engine := r.engine
	newGroup := &RouterGroup{
		prefix: r.prefix + prefix,
		parent: r,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (r *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := r.prefix + comp
	r.engine.router.addRoute(method, pattern, handler)
}

func (r *RouterGroup) GET(pattern string, handler HandlerFunc) {
	r.addRoute("GET", pattern, handler)
}

func (r *RouterGroup) POST(pattern string, handler HandlerFunc) {
	r.addRoute("POST", pattern, handler)
}

// Use is defined to add middleware to the group
func (r *RouterGroup) Use(middlewares ...HandlerFunc) {
	r.middlewares = append(r.middlewares, middlewares...)
}

// create static handler
func (r *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(r.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) {
		file := ctx.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

// Static serve static files
func (r *RouterGroup) Static(relativePath string, root string) {
	handler := r.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "*filepath")
	r.GET(urlPattern, handler)
}
