// Package cors/wrapper/gin provides gin.HandlerFunc to handle CORS related
// requests as a wrapper of github.com/rs/cors handler.
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

// Options is a configuration container to setup the CORS middleware.
type Options = cors.Options

// corsMiddleware is a wrapper of cors.Cors handler which preserves information
// about configured 'optionPassthrough' option.
type corsMiddleware struct {
	*cors.Cors
	optionPassthrough bool
}

// build transforms wrapped cors.Cors handler into Gin middleware.
func (c corsMiddleware) build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		if !c.optionPassthrough &&
			ctx.Request.Method == http.MethodOptions &&
			ctx.GetHeader("Access-Control-Request-Method") != "" {
			// Abort processing next Gin middlewares.
			ctx.AbortWithStatus(http.StatusOK)
		}
	}
}

// NewAllowAllCorsMiddleware creates a new CORS Gin middleware with permissive configuration
// allowing all origins with all standard methods with any header and
// credentials.
func NewAllowAllCorsMiddleware() gin.HandlerFunc {
	return corsMiddleware{Cors: cors.AllowAll()}.build()
}

// NewDefaultCorsMiddleware creates a new CORS Gin middleware with default options.
func NewDefaultCorsMiddleware() gin.HandlerFunc {
	return corsMiddleware{Cors: cors.Default()}.build()
}

// NewCorsMiddleware creates a new CORS Gin middleware with the provided options.
func NewCorsMiddleware(options Options) gin.HandlerFunc {
	return corsMiddleware{cors.New(options), options.OptionsPassthrough}.build()
}
