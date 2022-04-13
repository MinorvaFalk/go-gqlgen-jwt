package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MinorvaFalk/go-gqlgen-jwt/graph/model"
	"github.com/MinorvaFalk/go-gqlgen-jwt/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

const (
	ginContextKey = "GinContextKey"
)

type Middlewares struct {
	jwtAuth *utils.JwtAuth
}

type CustomClaims struct {
	*jwt.StandardClaims
	TokenType string
	model.User
}

func NewMiddlewares(jwtAuth *utils.JwtAuth) *Middlewares {
	return &Middlewares{jwtAuth: jwtAuth}
}

func (m *Middlewares) SignKey(username string) (string, error) {
	c := CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Minute).Unix(),
		},
		"level1",
		model.User{
			ID:   "1",
			Name: username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)

	return token.SignedString(m.jwtAuth.PrivateKey)
}

func (m *Middlewares) JwtMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if !strings.Contains(header, "Bearer") {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid Token Supplied")
	}

	tokenString := header[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

		if method, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("invalid token signature")
		} else if method != jwt.SigningMethodRS256 {
			return nil, errors.New("invalid token signature")
		}

		return m.jwtAuth.PublicKey, nil
	})

	if err != nil {
		c.AbortWithError(http.StatusForbidden, err)
	}

	if !token.Valid {
		c.AbortWithError(http.StatusForbidden, err)
	}

	c.Next()
}

func (m *Middlewares) GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
