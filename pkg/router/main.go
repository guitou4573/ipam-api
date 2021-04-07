package router

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guitou4573/ipam/pkg/controllers"
	"github.com/guitou4573/ipam/pkg/middlewares"
	"github.com/guitou4573/ipam/pkg/stores"
	"github.com/guitou4573/ipam/pkg/utils"
	"github.com/meatballhat/negroni-logrus"
	"github.com/urfave/negroni"
)

type Router struct {
	SubnetStore *stores.SubnetStore
  TokenStore *stores.TokenStore
	Router *gin.Engine
}

type RouterOptions struct {
	SubnetStore *stores.SubnetStore
  TokenStore *stores.TokenStore
}

func NewRouter(options RouterOptions) *Router{
	r := &Router{}
	r.SubnetStore = options.SubnetStore
  r.TokenStore = options.TokenStore

	r.Router = gin.New()
	r.Router.Use(gin.Recovery())

	r.initTracing()
	r.initCors()
	r.initLogs()

	midAuth := r.initAuth()
	networkController := controllers.NewNetworkController(
		controllers.NetworkControllerOptions{
			SubnetStore: r.SubnetStore,
		},
	)

  r.Router.NoRoute(func(c *gin.Context) {
    c.JSON(404, gin.H{
  		"error": "route not found",
  	})
  })
	r.initRoutes(midAuth, networkController)
	return r
}

func (r *Router) initTracing() {
	tracingMiddleware := middlewares.NewTracingMiddleware()
	r.Router.Use(tracingMiddleware)
}

func (r *Router) initCors() {
	// setup cors middleware
	corsMiddleware := middlewares.NewCorsMiddleware(middlewares.Options{
		AllowedOrigins: strings.Split(utils.EnvGet("CORS_ORIGIN", ""), ","),
		AllowedHeaders: []string{"*"},
		AllowCredentials: true,
		MaxAge: 600,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	r.Router.Use(corsMiddleware)
}

func (r *Router) initLogs() {
	// Set up a request logger, useful for debugging
	n := negroni.New()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(r.Router)
}

func (r *Router) initAuth() *middlewares.AuthMiddleware{
	midauth, err := middlewares.NewAuthMiddleware(r.TokenStore)
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return midauth
}

func (r *Router) initRoutes(midauth *middlewares.AuthMiddleware, networkController *controllers.NetworkController) {
	v1 := r.Router.Group("/v1")
	// v1.Use(midauth.Middleware.MiddlewareFunc())
	// {
    v1.GET("/vpc", networkController.ListVPCs)
		v1.GET("/vpc/:vpcid/subnet", networkController.ListSubnets)
		// v1.GET("/room/:roomid", messageController.ListMessages)
    // v1.GET("/room/:roomid/message", messageController.ListMessages)
		// v1.POST("/message", messageController.SendMessage)
	// }
}

func (r *Router) Start() {
	r.Router.RunTLS(":8000", "static/certs/cert.pem", "static/certs/key.pem")
}
