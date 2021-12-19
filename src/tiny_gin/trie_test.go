package tiny_gin

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Error("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/somePath")
	if n == nil {
		t.Error("nil shouldn't be returned")
	}

	if n.pattern != "/hello/:name" {
		t.Error("should match /hello/:name")
	}

	if ps["name"] != "somePath" {
		t.Error("name should be equal to 'somePath'")
	}
	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

	n, ps = r.getRoute("GET", "/assets/somepath/somefile")
	if ps["filepath"] != "somepath/somefile" {
		t.Error("filepath should be equal to 'somepath/somefile'")
	}
}
