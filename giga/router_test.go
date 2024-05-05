package giga

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
		t.Fatal("test parsePattern failed")
	}
}

func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	node, params := r.getRoute("GET", "/hello/makabaka")

	if node == nil {
		t.Fatal("nil shouldn't be returned")
	}

	if node.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	if params["name"] != "makabaka" {
		t.Fatal("name should be equal to 'makabaka'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", node.pattern, params["name"])
}
