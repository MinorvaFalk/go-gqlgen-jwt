package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"fmt"

	"github.com/MinorvaFalk/go-gqlgen-jwt/graph/model"
	"github.com/MinorvaFalk/go-gqlgen-jwt/middleware"
	"github.com/gin-gonic/gin"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	todos      []*model.Todo
	Middleware *middleware.Middlewares
}

func (r *Resolver) GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}
