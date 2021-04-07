package middlewares

import (
	"log"
	"time"

	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/guitou4573/ipam/pkg/models"
	"github.com/guitou4573/ipam/pkg/stores"
)

type AuthMiddleware struct {
	Middleware *jwt.GinJWTMiddleware
	TokenStore *stores.TokenStore
}

type JWTLogin struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "sub"
var realm = ""
var secretKey = []byte("00000000")

func NewAuthMiddleware(tokenStore *stores.TokenStore) (*AuthMiddleware, error) {
	midAuth := &AuthMiddleware{}
	midAuth.TokenStore = tokenStore

	// the jwt middleware
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: realm,
		Key: secretKey,
		SigningAlgorithm: "HS512",
		Timeout: time.Hour,
		MaxRefresh: time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) (jwt.MapClaims) {return midAuth.PayloadFunc(data)},
		IdentityHandler:	midAuth.IdentityHandler,
		Authenticator: func(c *gin.Context) (interface{}, error) {return midAuth.Authenticator(c)},
		Authorizator: midAuth.Authorizator,
		Unauthorized: midAuth.Unauthorized,
		TokenLookup: "header: Authorization", //header: Authorization, query: token, cookie: jwt
		TokenHeadName: "Bearer",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	midAuth.Middleware = middleware
	if err != nil {
		return nil, err
	}

	return midAuth, nil
}

func (a *AuthMiddleware) Unauthorized (c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

func (a *AuthMiddleware) Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*models.Token); ok {
    _, err := a.TokenStore.Get(v.GetToken())
		if(err != nil) {
			return false
		}
		return true
	}
	return false
}

func (a *AuthMiddleware) Authenticator(c *gin.Context) (interface{}, error) {
	// disabled, using only auth app for login
	var loginVals JWTLogin
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	return nil, jwt.ErrFailedAuthentication
}

func (a *AuthMiddleware) IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	if (claims[identityKey] != nil ) {
		// load user from DB
		_, err := a.TokenStore.Get(claims[identityKey].(string))
		if(err != nil) {
			log.Println(err)
			return nil
		}
		return nil
	}
	log.Println("Identity not provided in token")
	return nil
}

func (a *AuthMiddleware) PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.Token); ok {
		return jwt.MapClaims{
			identityKey: v.Token,
		}
	}
	return jwt.MapClaims{}
}
