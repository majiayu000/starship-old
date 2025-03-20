package middleware

import (
	"github.com/gin-gonic/gin"
)

// Middleware represents a middleware function for Gin
type Middleware func(gin.HandlerFunc) gin.HandlerFunc

// Chain applies multiple middlewares to a handler
func Chain(handler gin.HandlerFunc, middlewares ...Middleware) gin.HandlerFunc {
	for _, middleware := range middlewares {
		if middleware != nil {
			handler = middleware(handler)
		}
	}
	return handler
}
