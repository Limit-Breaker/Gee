package greet

import (
	"giga"
	"net/http"
)

type HandlerGreet struct {
}

func (h *HandlerGreet) Hello(c *giga.Context) {
	// expect /hello/makabaka
	c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
}
