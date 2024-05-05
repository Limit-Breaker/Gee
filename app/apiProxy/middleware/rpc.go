package middleware

import (
	"giga"
)

// MiddlewareRpc 将rpc client实例存在Keys中
func MiddlewareRpc(services map[string]interface{}) giga.HandlerFunc {
	return func(c *giga.Context) {
		// 将rpc client实例存在Keys中
		c.Keys = make(map[string]interface{})
		for k, v := range services {
			c.Keys[k] = v
		}
		c.Next()
	}
}
