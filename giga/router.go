package giga

import (
	"fmt"
	"net/http"
	"strings"
)

// roots key eg, roots['GET'] roots['POST']
// handlers key eg, ['GET-/index/:id/detail', 'POST-/user/login']
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node, 0),
		handlers: make(map[string]HandlerFunc, 0),
	}
}

// Only one * is allowed
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 通配*。例如 /path/*filepath, 只能允许一个*作为最后一个节点
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	if _, ok := r.roots[method]; !ok {
		// 如果不存在该方法的根节点
		r.roots[method] = &node{}
	}

	// 将路由插入
	r.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 解析了:和*两种匹配符的参数，返回一个 map 。
// 例如/index/1/detail匹配到/index/:id/detail，返回的解析结果为：{id: "1"}
// 例如path/file/log 匹配到/path/*filepath，解析结果为{filepath: "file/log"}
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string, 0)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0)
	if node == nil {
		return nil, nil
	}
	parts := parsePattern(node.pattern)
	for index, part := range parts {
		// 例如/index/1/detail匹配到/index/:id/detail，返回的解析结果为：{id: "1"}
		if part[0] == ':' {
			params[part[1:]] = searchParts[index]
		}
		// 例如path/file/log 匹配到/path/*filepath，解析结果为{filepath: "file/log"}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}
	return node, params
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
		// 找到请求处理函数
		key := c.Method + "-" + node.pattern
		fmt.Printf("r.handers: %+v \n", r.handlers)
		// 将请求处理函数也加入到context.handler中
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// 执行中间件链路函数和请求处理函数
	c.Next()
}
